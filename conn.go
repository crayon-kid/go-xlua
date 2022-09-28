package lua

import (
	"skywalk/shared/protoc"
	"skywalk/shared/srpc"
)

//返回命令中的错误
type Error string

func (err Error) Error() string { return string(err) }

type Conn interface {
	//关闭服务
	Close() error
	//战斗验证接口
	DoCall(commandName string, args ...interface{}) (reply interface{}, err error)
}
