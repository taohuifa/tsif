package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	// "runtime/debug"
	"runtime"

	"bytes"
)

// 日志环境
type LogContext struct {
	logLevel int // 日志等级

	useLogger bool   // 使用日志输出
	logPath   string // 日志路径
	logName   string // 日志名

	loggers map[int]*log.Logger // 日志器
}

const (
	LOG_DEBUG = 0 // debug
	LOG_INFO  = 1 // info
	LOG_WARN  = 2 // warn
	LOG_ERROR = 3 // error
)

// 获取log名称
func getLogLvName(level int) string {
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

// 检测日志等级
func checkLogLv(loglv int) bool {
	switch loglv {
	case LOG_DEBUG:
		return true
	case LOG_INFO:
		return true
	case LOG_WARN:
		return true
	case LOG_ERROR:
		return true
	}
	return false
}

// 获取堆栈信息
func getStack(frame int, depth int) string {
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

// 初始化日志
func createLogger(loglv int, logPath string, logName string) *log.Logger {
	//check log path
	if logPath == "" {
		panic(errors.New("empty log path " + logPath))

	}
	// log file
	logLvName := getLogLvName(loglv)
	logfile := fmt.Sprintf("%s/%s_%s_%s.txt", logPath, logName, strings.ToLower(logLvName), time.Now().Local().Format("2006010215"))

	// check path
	_, errPath := os.Stat(logPath)
	if errPath != nil {
		// 错误, 文件夹不存在
		errMkdir := os.Mkdir(logPath, os.ModePerm)
		if errMkdir != nil {
			panic(errMkdir)
		}
	}

	// create log file
	logFile, err := os.Create(logfile)
	if err != nil {
		panic(errors.New("error log file! loglv=" + fmt.Sprint(loglv) + " file=" + logfile + "."))
	}
	// l := log.New(logFile, "", log.Llongfile|log.Ldate|log.Ltime|log.Lmicroseconds)
	l := log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	return l
}

func From(loglv int, logPath string, logName string) *LogContext {
	if !checkLogLv(loglv) {
		panic(errors.New("error log level " + fmt.Sprint(loglv)))
	}

	// create
	logObj := LogContext{logLevel: loglv, logPath: logPath, logName: logName}
	logObj.loggers = make(map[int]*log.Logger)
	// check use logger
	if logPath != "" && logName != "" {
		logObj.useLogger = true
		// test by init logger
		//logObj.GetLogger(LOG_DEBUG)
	}
	return &logObj
}

func (this *LogContext) GetLogger(loglv int) *log.Logger {
	// check
	if !checkLogLv(loglv) {
		panic(errors.New("error log level " + fmt.Sprint(loglv)))
	}
	// get logger
	l, ok := this.loggers[loglv]
	if ok {
		return l
	}
	// create init
	l = createLogger(loglv, this.logPath, this.logName)
	this.loggers[loglv] = l
	return l
}

// 输出函数
func (this *LogContext) Write(loglv int, frame int, msg string, ext string) {
	//check log level
	if loglv < this.logLevel {
		return
	}

	var buffer bytes.Buffer
	buffer.WriteString(getLogLvName(loglv))
	buffer.WriteString(" ")
	buffer.WriteString(time.Now().Local().Format("2006-01-02 15:04:05"))
	buffer.WriteString(" ")
	buffer.WriteString(getStack(frame, 1))
	buffer.WriteString(":")
	buffer.WriteString("\r\n")
	buffer.WriteString(msg)
	if ext != "" {
		buffer.WriteString("\r\n\t")
		// buffer.WriteString(get_stack(frame+1, 99))
		buffer.WriteString(ext)
	}
	// log.Printf(buffer.String())
	fmt.Println(buffer.String())
	if this.useLogger {
		// 遍历写入日志
		for i := loglv; i >= this.logLevel; i-- {
			this.GetLogger(i).Print(buffer.String())
		}
	}

	// debug.PrintStack()
}

// 输出函数和堆栈
func (this *LogContext) WriteStack(loglv int, frame int, msg string) {
	this.Write(loglv, frame, msg, getStack(frame+1, 99))
}

func (this *LogContext) Debug(msg string) {
	this.Write(LOG_DEBUG, 4, msg, "")
}
func (this *LogContext) Info(msg string) {
	this.Write(LOG_INFO, 4, msg, "")
}
func (this *LogContext) Warn(msg string) {
	this.Write(LOG_WARN, 4, msg, "")
}
func (this *LogContext) Error(msg string) {
	this.Write(LOG_ERROR, 4, msg, "")
}

func (this *LogContext) DebugStack(msg string) {
	this.WriteStack(LOG_DEBUG, 4, msg)
}
func (this *LogContext) InfoStack(msg string) {
	this.WriteStack(LOG_INFO, 4, msg)
}
func (this *LogContext) WarnStack(msg string) {
	this.WriteStack(LOG_WARN, 4, msg)
}
func (this *LogContext) ErrorStack(msg string) {
	this.WriteStack(LOG_ERROR, 4, msg)
}
