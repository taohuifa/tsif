package log

// var Instance LogContext
// var Log LogContext = *From(LOG_DEBUG, "./", "log")
var Log LogContext = *From(LOG_INFO, "", "log")

func Debug(msg string) {
	Log.Debug(msg)
}
func Info(msg string) {
	Log.Info(msg)
}
func Warn(msg string) {
	Log.Warn(msg)
}
func Error(msg string) {
	Log.Error(msg)
}

func Write(loglv int, frame int, msg string) {
	Log.Write(loglv, frame, msg, "")
}

func DebugStack(msg string) {
	Log.DebugStack(msg)
}
func InfoStack(msg string) {
	Log.InfoStack(msg)
}
func WarnStack(msg string) {
	Log.WarnStack(msg)
}
func ErrorStack(msg string) {
	Log.ErrorStack(msg)
}
func WriteStack(loglv int, frame int, msg string) {
	Log.WriteStack(loglv, frame, msg)
}
