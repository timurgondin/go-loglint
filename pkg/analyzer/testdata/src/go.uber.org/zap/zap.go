package zap

type Logger struct{}

func (l *Logger) Info(msg string, args ...interface{})  {}
func (l *Logger) Error(msg string, args ...interface{}) {}
func (l *Logger) Warn(msg string, args ...interface{})  {}
func (l *Logger) Debug(msg string, args ...interface{}) {}

func L() *Logger {
	return &Logger{}
}

func Info(msg string, args ...interface{})  {}
func Error(msg string, args ...interface{}) {}
func Warn(msg string, args ...interface{})  {}
func Debug(msg string, args ...interface{}) {}
