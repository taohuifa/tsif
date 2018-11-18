package log

// var Instance LogContext
// var Log LogContext = *From(LOG_DEBUG, "./", "log")
var Log LogContext = *From(LOG_INFO, "", "log")

func Debug(params ...interface{}) {
	Log.Debug(params...)
}
func Info(params ...interface{}) {
	Log.Info(params...)
}
func Warn(params ...interface{}) {
	Log.Warn(params...)
}
func Error(params ...interface{}) {
	Log.Error(params...)
}

func Write(loglv int, frame int, params ...interface{}) {
	Log.Write(loglv, frame, "", params...)
}

func Debugf(format string, params ...interface{}) {
	Log.Debugf(format, params...)
}
func Infof(format string, params ...interface{}) {
	Log.Infof(format, params...)
}
func Warnf(format string, params ...interface{}) {
	Log.Warnf(format, params...)
}
func Errorf(format string, params ...interface{}) {
	Log.Errorf(format, params...)
}

func Writef(loglv int, frame int, format string, params ...interface{}) {
	Log.Writef(loglv, frame+1, "", format, params...)
}

func Debugs(params ...interface{}) {
	Log.Debugs(params...)
}
func Infos(params ...interface{}) {
	Log.Infos(params...)
}
func Warns(params ...interface{}) {
	Log.Warns(params...)
}
func Errors(params ...interface{}) {
	Log.Errors(params...)
}
func Writes(loglv int, frame int, params ...interface{}) {
	Log.Writes(loglv, frame, params...)
}
