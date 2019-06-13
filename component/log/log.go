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

// 日志对象
type logItem struct {
	logger      *log.Logger // 日志器
	logFile     *os.File    // 日志文件
	logFileName string      // 当前日志文件名
	checkTime   int64       // 上次检测时间
}

// 日志环境
type LogContext struct {
	logLevel int // 日志等级

	useLogger bool   // 使用日志输出
	logPath   string // 日志路径
	logName   string // 日志名

	loggers map[int]*logItem // 日志器
}

const (
	LOG_DEBUG = 0 // debug
	LOG_INFO  = 1 // info
	LOG_WARN  = 2 // warn
	LOG_ERROR = 3 // error

	LOG_SHOW_FUNC  = false // 日志是否带文件输出
	LOG_BASE_FRAME = 5     // 基础frame
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
		// 输出带文件
		if LOG_SHOW_FUNC {
			buffer.WriteString(fmt.Sprintf("at %s(%s:%d)", runtime.FuncForPC(pc).Name(), file, line))
		} else {
			buffer.WriteString(fmt.Sprintf("at %s:%d", file, line))
		}
	}
	return buffer.String()
}

// 初始化日志
func (this *logItem) init(loglv int, logFileName string) error {
	//check log path
	if logFileName == "" {
		panic(errors.New("empty log path " + logFileName))
	}
	this.logFileName = logFileName

	// create log file
	// logFile, err := os.Create(this.logFileName)
	logFile, err := os.OpenFile(this.logFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666) // 打开追加文件
	if err != nil {
		return errors.New("error log file! loglv=" + fmt.Sprint(loglv) + " file=" + this.logFileName + ".")
	}
	// this.logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	this.logger = log.New(logFile, "", 0)
	this.logFile = logFile
	this.checkTime = time.Now().Unix()
	return nil
}

// 销毁
func (this *logItem) destory() error {
	// this.logger.close()
	// this.logger.GetFile()
	this.logger = nil
	this.logFileName = ""
	this.checkTime = 0
	this.logFile.Close()
	this.logFile = nil
	return nil
}

// to string
func ToString(params ...interface{}) string {
	pnum := len(params)
	if pnum <= 0 {
		return ""
	} else if pnum == 1 {
		return fmt.Sprint(params[0])
	}
	// write string
	var buffer bytes.Buffer
	for i := 0; i < pnum; i++ {
		if i > 0 {
			// buffer.WriteString("\t")
			buffer.WriteString(" ")
		}
		str := fmt.Sprint(params[i])
		buffer.WriteString(str)
	}
	return buffer.String()
}

func From(loglv int, logPath string, logName string) *LogContext {
	if !checkLogLv(loglv) {
		panic(errors.New("error log level " + fmt.Sprint(loglv)))
	}

	// create
	logObj := LogContext{logLevel: loglv, logPath: logPath, logName: logName}
	logObj.loggers = make(map[int]*logItem)
	// check use logger
	if logPath != "" && logName != "" {
		logObj.useLogger = true
		// test by init logger
		//logObj.GetLogger(LOG_DEBUG)
	}
	return &logObj
}

// get log file name
func getLogFileName(loglv int, logPath string, logName string) string {
	logLvName := getLogLvName(loglv)
	logFileName := fmt.Sprintf("%s/%s_%s_%s.txt", logPath, logName, strings.ToLower(logLvName), time.Now().Local().Format("20060102")) // 天
	// logFileName := fmt.Sprintf("%s/%s_%s_%s.txt", logPath, logName, strings.ToLower(logLvName), time.Now().Local().Format("2006010215")) // 小时
	// logFileName := fmt.Sprintf("%s/%s_%s_%s.txt", logPath, logName, strings.ToLower(logLvName), time.Now().Local().Format("200601021504")) // 分钟
	return logFileName
}

func (this *LogContext) GetLogger(loglv int) *log.Logger {
	// check
	if !checkLogLv(loglv) {
		panic(errors.New("error log level " + fmt.Sprint(loglv)))
	}

	// get logger
	l, ok := this.loggers[loglv]
	if ok {
		// 检测时间, 60s检测一次
		nowTime := time.Now().Unix()
		if nowTime-l.checkTime < 10 {
			return l.logger // 未到检测时间
		}
		l.checkTime = nowTime

		// 检测日志是否对应(检测日志更新)
		logFileName := getLogFileName(loglv, this.logPath, this.logName)
		if l.logFileName == logFileName {
			return l.logger // 文件一致, 可直接用
		}
		// 日志文件失效, 重新创建1个
		l.destory()
		this.loggers[loglv] = nil
		l = nil
	}

	// check path
	_, errPath := os.Stat(this.logPath)
	if errPath != nil {
		// 错误, 文件夹不存在
		errMkdir := os.Mkdir(this.logPath, os.ModePerm)
		if errMkdir != nil {
			panic(errMkdir)
			// return nil
		}
	}

	// create log
	l = &logItem{}
	if l == nil {
		panic(errors.New("create log fail, level " + fmt.Sprint(loglv)))
		// return nil
	}
	logFileName := getLogFileName(loglv, this.logPath, this.logName)

	// init log item
	lerr := l.init(loglv, logFileName)
	if lerr != nil {
		panic(lerr)
		// return nil
	}
	// 获取日志器
	this.loggers[loglv] = l
	return l.logger
}

// 输出函数
func (this *LogContext) write(loglv int, frame int, ext string, msg string) {
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
	// buffer.WriteString("\r\n")	// 换行, 搜索起来不带时间不好
	buffer.WriteString(" ")
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

// 输出函数
func (this *LogContext) Write(loglv int, frame int, ext string, params ...interface{}) {
	this.write(loglv, frame, ext, ToString(params...))
}

// 输出函数
func (this *LogContext) Writef(loglv int, frame int, ext string, format string, params ...interface{}) {
	this.write(loglv, frame, ext, fmt.Sprintf(format, params...))
}

// 输出函数和堆栈
func (this *LogContext) Writes(loglv int, frame int, params ...interface{}) {
	this.write(loglv, frame, getStack(frame+1, 99), ToString(params...))
}

// 输出函数和堆栈
func (this *LogContext) Writesf(loglv int, frame int, format string, params ...interface{}) {
	this.Write(loglv, frame, getStack(frame+1, 99), fmt.Sprintf(format, params...))
}

func (this *LogContext) Debug(params ...interface{}) {
	this.Write(LOG_DEBUG, LOG_BASE_FRAME, "", params...)
}
func (this *LogContext) Info(params ...interface{}) {
	this.Write(LOG_INFO, LOG_BASE_FRAME, "", params...)
}
func (this *LogContext) Warn(params ...interface{}) {
	this.Write(LOG_WARN, LOG_BASE_FRAME, "", params...)
}
func (this *LogContext) Error(params ...interface{}) {
	this.Write(LOG_ERROR, LOG_BASE_FRAME, "", params...)
}

func (this *LogContext) Debugf(format string, params ...interface{}) {
	this.Writef(LOG_DEBUG, LOG_BASE_FRAME, "", format, params...)
}
func (this *LogContext) Infof(format string, params ...interface{}) {
	this.Writef(LOG_INFO, LOG_BASE_FRAME, "", format, params...)
}
func (this *LogContext) Warnf(format string, params ...interface{}) {
	this.Writef(LOG_WARN, LOG_BASE_FRAME, "", format, params...)
}
func (this *LogContext) Errorf(format string, params ...interface{}) {
	this.Writef(LOG_ERROR, LOG_BASE_FRAME, "", format, params...)
}

func (this *LogContext) Debugs(params ...interface{}) {
	this.Writes(LOG_DEBUG, LOG_BASE_FRAME, params...)
}
func (this *LogContext) Infos(params ...interface{}) {
	this.Writes(LOG_INFO, LOG_BASE_FRAME, params...)
}
func (this *LogContext) Warns(params ...interface{}) {
	this.Writes(LOG_WARN, LOG_BASE_FRAME, params...)
}
func (this *LogContext) Errors(params ...interface{}) {
	this.Writes(LOG_ERROR, LOG_BASE_FRAME, params...)
}
