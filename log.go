package lua

import (
	"fmt"
)

type LogType int8

const (
	LogType_Info = iota + 1
	LogType_Warn
	LogType_Error
)

func LogInfo(l *State) int {
	return Log(l, LogType_Info)
}

func LogWarning(l *State) int {
	return Log(l, LogType_Warn)
}

func LogError(l *State) int {
	return Log(l, LogType_Error)
}

func Log(l *State, logType LogType) int {
	return 0
	n := l.GetTop()
	l.GetGlobal("tostring")
	var s string
	for i := 1; i <= n; i++ {
		l.PushValue(-1)
		l.PushValue(i)
		err := l.Call(1, 1)
		if err != nil {
			fmt.Println(err.Error())
			return 1
		}
		s += l.ToString(-1)
		if i != n {
			s += "\t"
		}
		l.Pop(1)
	}
	switch logType {
	case LogType_Info:
		fmt.Printf("\033[1;32m %s:%s \033[0m\n", "Info", s)
	case LogType_Warn:
		fmt.Printf("\033[1;33m %s:%s \033[0m\n", "Warn", s)
	case LogType_Error:
		fmt.Printf("\033[1;31m %s:%s \033[0m\n", "Error", s)
	default:
		fmt.Printf("\u001B[1;35m %s:%s \u001B[0m\n", "UnDefine", s)
	}
	return 0
}
