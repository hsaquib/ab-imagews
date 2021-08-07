package rest_log

import "sync"

type LogLevel string

var once sync.Once

const (
	Info  LogLevel = "Info"
	Warn  LogLevel = "Warn"
	Debug LogLevel = "Debug"
	Error LogLevel = "Error"
	Fatal LogLevel = "Fatal"
)

type Logger interface {
	Info(fn, tid string, msg string)
	InfoPretty(fn, tid string, msg string)
	Warn(fn, tid string, msg string)
	WarnPretty(fn, tid string, msg string)
	Error(fn, tid string, msg string)
	ErrorPretty(fn, tid string, msg string)
	Print(level LogLevel, fn, tid string, msg string)
}

var lgr Logger
var applicationName = ""
var isVerbose bool

func New(verbose bool, appName string) Logger {
	logger := NewZeroLevelLogger(verbose, appName)
	logger.Info("GetDefaultStructLogger", "init", "Running in verbose mode")
	return logger
}

func Init(verbose bool, appName string) {
	applicationName = appName
	isVerbose = verbose
}

func GetLogger() Logger {
	once.Do(func() {
		if applicationName == "" {
			isVerbose = true
			applicationName = "DefaultStructLogger"
		}
		lgr = NewZeroLevelLogger(isVerbose, applicationName)
		if isVerbose {
			lgr.Info("Logger", "init", "Running in verbose mode")
		} else {
			lgr.Info("Logger", "init", "Running in non-verbose mode")
		}
	})
	return lgr
}
