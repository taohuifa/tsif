package log

import (
	"fmt"
	"log"

	// "runtime/debug"
	"runtime"

	"bytes"
)

// 日志环境
type LogContext struct {
	LogLevel int // 日志等级
}

const (
	LOG_DEBUG = 0 // debug
	LOG_INFO  = 1 // info
	LOG_WARN  = 2 // warn
	LOG_ERROR = 3 // error
)

// 获取log名称
func get_logname(level int) string {
	switch level {
	case LOG_DEBUG:
		return "DEBUG"
	case LOG_INFO:
		return "INFO"
	case LOG_WARN:
		return "WARN"
	case LOG_ERROR:
		return "ERROR"
	}
	return "UNKNOW"
}

// 获取堆栈信息
func get_stack(frame int, depth int) string {
	var buffer bytes.Buffer //Buffer是一个实现了读写方法的可变大小的字节缓冲

	// szHead = fmt.Sprintf("%s(%s:%d)", runtime.FuncForPC(pc).Name(), file, line)
	for i := 0; i < depth; i++ {
		pc, file, line, ok := runtime.Caller(i + frame)
		if !ok {
			break
		}
		if i > 0 {
			buffer.WriteString("\r\n\t")
		}
		buffer.WriteString(fmt.Sprintf("at %s(%s:%d)", runtime.FuncForPC(pc).Name(), file, line))
	}
	return buffer.String()
}

// 输出函数
func (this *LogContext) Write(level int, msg string, showstack bool) {
	var buffer bytes.Buffer

	buffer.WriteString(get_logname(level))
	buffer.WriteString(" ")
	buffer.WriteString(get_stack(4, 1))
	buffer.WriteString(":")
	buffer.WriteString("\r\n")
	buffer.WriteString(msg)
	if showstack {
		buffer.WriteString("\r\n\t")
		buffer.WriteString(get_stack(5, 99))
	}
	log.Printf(buffer.String())
	// debug.PrintStack()
}

func (this *LogContext) Debug(msg string, showstack bool) {
	this.Write(LOG_DEBUG, msg, showstack)
}

func (this *LogContext) Info(msg string, showstack bool) {
	this.Write(LOG_INFO, msg, showstack)
}
func (this *LogContext) Warn(msg string, showstack bool) {
	this.Write(LOG_WARN, msg, showstack)
}
func (this *LogContext) Error(msg string, showstack bool) {
	this.Write(LOG_ERROR, msg, showstack)
}

var Log LogContext

func Debug(msg string, showstack bool) {
	Log.Debug(msg, showstack)
}
func Info(msg string, showstack bool) {
	Log.Info(msg, showstack)
}
func Warn(msg string, showstack bool) {
	Log.Warn(msg, showstack)
}
func Error(msg string, showstack bool) {
	Log.Error(msg, showstack)
}

// func Debug(msg string) {
// 	Instance.Debug(msg, false)
// }
// func Info(msg string) {
// 	Instance.Info(msg, false)
// }
// func Warn(msg string) {
// 	Instance.Warn(msg, false)
// }
// func Error(msg string) {
// 	Instance.Error(msg, false)
// }
