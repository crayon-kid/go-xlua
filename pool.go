package lua

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var nowFunc = time.Now // for testing

//超过最大活跃数量时
var ErrPoolExhausted = errors.New("golua: state pool exhausted")

//lua环境被关闭
var (
	errConnClosed = errors.New("golua: state closed")
)

type Pool struct {
	// 最大等待连接数
	MaxIdle int
	// 最大活跃连接数 （=0无上限）
	MaxActive int
	// 等待连接的过期时间(=0不过期)
	//IdleTimeout time.Duration
	//是否等待
	Wait bool

	mu           sync.Mutex
	closed       bool
	active       int
	initOnce     sync.Once
	ch           chan struct{}
	idle         idleList
	waitCount    int64
	waitDuration time.Duration
}

//生成池子
func NewPool(maxConn int) *Pool {
	return &Pool{MaxIdle: maxConn, MaxActive: maxConn, Wait: true}
}

//预热
func (p *Pool) WarmUp() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := 0; i < p.MaxIdle; i++ {
		l, err := p.create()
		if err != nil {
			return err
		}
		p.active++
		conn := &poolConn{l: l}
		p.idle.pushFront(conn)
	}
	return nil
}

//获得lua连接
func (p *Pool) Get() Conn {
	//排队以及信息记录
	waited, err := p.waitVacantConn()
	if err != nil {
		return errorConn{err}
	}

	p.mu.Lock()
	if waited > 0 {
		p.waitCount++
		p.waitDuration += waited
	}

	//if p.IdleTimeout > 0 {
	//	//清除长时间处于等待的连接
	//	n := p.idle.count
	//	for i := 0; i < n && p.idle.back != nil && p.idle.back.t.Add(p.IdleTimeout).Before(nowFunc()); i++ {
	//		pc := p.idle.back
	//		p.idle.popBack()
	//		p.mu.Unlock()
	//		pc.l.Close()
	//		p.mu.Lock()
	//		p.active--
	//	}
	//}

	//从队首获取可用lua连接
	for p.idle.front != nil {
		pc := p.idle.front
		p.idle.popFront()
		p.mu.Unlock()
		return &activeConn{p: p, pc: pc}
	}

	//池子已关闭
	if p.closed {
		p.mu.Unlock()
		err := errors.New("golua: get on closed pool")
		return errorConn{err}
	}

	//超过最大活跃上限
	if !p.Wait && p.MaxActive > 0 && p.active >= p.MaxActive {
		p.mu.Unlock()
		return errorConn{ErrPoolExhausted}
	}

	//新生成lua连接
	p.active++
	p.mu.Unlock()
	l, err := p.create()
	if err != nil {
		p.mu.Lock()
		p.active--
		if p.ch != nil && !p.closed {
			p.ch <- struct{}{}
		}
		p.mu.Unlock()
		return errorConn{err}
	}
	return &activeConn{p: p, pc: &poolConn{l: l}}
}

//池子使用状况
type PoolStats struct {
	ActiveCount  int
	IdleCount    int
	WaitCount    int64
	WaitDuration time.Duration
}

func (p *Pool) Stats() PoolStats {
	p.mu.Lock()
	stats := PoolStats{
		ActiveCount:  p.active,
		IdleCount:    p.idle.count,
		WaitCount:    p.waitCount,
		WaitDuration: p.waitDuration,
	}
	p.mu.Unlock()

	return stats
}

//获取池子中拥有的连接数
func (p *Pool) ActiveCount() int {
	p.mu.Lock()
	active := p.active
	p.mu.Unlock()
	return active
}

//获取池子中处于可用状态的连接数
func (p *Pool) IdleCount() int {
	p.mu.Lock()
	idle := p.idle.count
	p.mu.Unlock()
	return idle
}

//关闭池子
func (p *Pool) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	p.active -= p.idle.count
	//pc := p.idle.front
	p.idle.count = 0
	p.idle.front, p.idle.back = nil, nil
	if p.ch != nil {
		close(p.ch)
	}
	p.mu.Unlock()
	//耗时很长 宿主程序关闭资源自然释放所以直接不调用了
	//for ; pc != nil; pc = pc.next {
	//	pc.l.Close()
	//}
	return nil
}

func (p *Pool) lazyInit() {
	p.initOnce.Do(func() {
		p.ch = make(chan struct{}, p.MaxActive)
		if p.closed {
			close(p.ch)
		} else {
			for i := 0; i < p.MaxActive; i++ {
				p.ch <- struct{}{}
			}
		}
	})
}

//等待排队
func (p *Pool) waitVacantConn() (waited time.Duration, err error) {
	if !p.Wait || p.MaxActive <= 0 {
		return 0, nil
	}

	p.lazyInit()

	wait := len(p.ch) == 0
	var start time.Time
	if wait {
		start = time.Now()
	}

	<-p.ch

	if wait {
		return time.Since(start), nil
	}
	return 0, nil
}

func (p *Pool) create() (*State, error) {
	l := NewState()
	l.OpenLibs()

	if err := p.load(l); err != nil {
		return nil, err
	}
	return l, nil
}

func (p *Pool) load(l *State) error {
	//添加lua代码搜索路径
	l.GetGlobal("package")
	l.GetField(-1, "path")
	oldPath := l.ToString(-1)
	l.PushString(fmt.Sprintf("%v;%v", oldPath, "luascripts/?.lua"))
	l.SetField(-3, "path")
	l.Pop(2)
	//引入第三方包
	l.RegistryDllFunc()
	//l.Register("_loginfo", LogInfo)
	//l.Register("_logwarning", LogWarning)
	//l.Register("_logerror", LogError)
	return nil
}

func (p *Pool) put(pc *poolConn, forceClose bool) error {
	p.mu.Lock()
	if !p.closed && !forceClose {
		//pc.t = nowFunc()
		p.idle.pushFront(pc)
		if p.idle.count > p.MaxIdle {
			pc = p.idle.back
			p.idle.popBack()
		} else {
			pc = nil
		}
	}

	if pc != nil {
		p.mu.Unlock()
		pc.l.Close()
		p.mu.Lock()
		p.active--
	}

	if p.ch != nil && !p.closed {
		p.ch <- struct{}{}
	}
	p.mu.Unlock()
	return nil
}

type activeConn struct {
	p  *Pool
	pc *poolConn
}

func (ac *activeConn) firstError(errs ...error) error {
	for _, err := range errs[:len(errs)-1] {
		if err != nil {
			return err
		}
	}
	return errs[len(errs)-1]
}

func (ac *activeConn) Close() error {
	pc := ac.pc
	if pc == nil {
		return nil
	}
	ac.pc = nil

	return ac.p.put(pc, false)
}

func (ac *activeConn) DoCall(commandName string, args ...interface{}) (reply interface{}, err error) {
	pc := ac.pc
	if pc == nil || pc.l == nil {
		return nil, errConnClosed
	}
	// do call
	return nil, nil
}

type errorConn struct{ err error }

func (ec errorConn) DoCall(commandName string, args ...interface{}) (reply interface{}, err error)
	return nil, ec.err
}

func (ec errorConn) Close() error { return nil }

type idleList struct {
	count       int
	front, back *poolConn
}

type poolConn struct {
	l *State
	//t          time.Time
	next, prev *poolConn
}

func (l *idleList) pushFront(pc *poolConn) {
	pc.next = l.front
	pc.prev = nil
	if l.count == 0 {
		l.back = pc
	} else {
		l.front.prev = pc
	}
	l.front = pc
	l.count++
}

func (l *idleList) popFront() {
	pc := l.front
	l.count--
	if l.count == 0 {
		l.front, l.back = nil, nil
	} else {
		pc.next.prev = nil
		l.front = pc.next
	}
	pc.next, pc.prev = nil, nil
}

func (l *idleList) popBack() {
	pc := l.back
	l.count--
	if l.count == 0 {
		l.front, l.back = nil, nil
	} else {
		pc.prev.next = nil
		l.back = pc.prev
	}
	pc.next, pc.prev = nil, nil
}
