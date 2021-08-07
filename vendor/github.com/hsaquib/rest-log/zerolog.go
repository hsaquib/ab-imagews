package rest_log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

type zeroLevelLogger struct {
	StdLog       zerolog.Logger
	ErrLog       zerolog.Logger
	PrettyStdLog zerolog.Logger
	PrettyErrLog zerolog.Logger
	verbose      bool
}

func NewZeroLevelLogger(verbose bool, appName string) Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	stdLog := zerolog.New(os.Stdout).With().CallerWithSkipFrameCount(3).Str("service", appName).Logger()
	errLog := zerolog.New(os.Stderr).With().CallerWithSkipFrameCount(3).Stack().Str("service", appName).Logger()
	prettyStdLog := log.With().CallerWithSkipFrameCount(3).Str("service", appName).Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})
	prettyErrorLog := log.With().CallerWithSkipFrameCount(3).Str("service", appName).Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})

	return &zeroLevelLogger{
		StdLog:       stdLog,
		ErrLog:       errLog,
		PrettyStdLog: prettyStdLog,
		PrettyErrLog: prettyErrorLog,
		verbose:      verbose,
	}
}

func (l zeroLevelLogger) Info(fn, tid string, msg string) {
	level := zerolog.Disabled
	if l.verbose {
		level = zerolog.InfoLevel
	}
	l.StdLog.WithLevel(level).Str("function", fn).Str("tid", tid).Msg(msg)
}

func (l zeroLevelLogger) InfoPretty(fn, tid string, msg string) {
	level := zerolog.Disabled
	if l.verbose {
		level = zerolog.InfoLevel
	}
	l.PrettyStdLog.WithLevel(level).Str("function", fn).Str("tid", tid).Msg(msg)
}

func (l zeroLevelLogger) Warn(fn, tid string, msg string) {
	level := zerolog.Disabled
	if l.verbose {
		level = zerolog.WarnLevel
	}
	l.StdLog.WithLevel(level).Str("function", fn).Str("tid", tid).Msg(msg)
}

func (l zeroLevelLogger) WarnPretty(fn, tid string, msg string) {
	level := zerolog.Disabled
	if l.verbose {
		level = zerolog.WarnLevel
	}
	l.PrettyStdLog.WithLevel(level).Str("function", fn).Str("tid", tid).Msg(msg)
}

func (l zeroLevelLogger) Error(fn, tid string, msg string) {
	level := zerolog.Disabled
	if l.verbose {
		level = zerolog.ErrorLevel
	}
	l.ErrLog.WithLevel(level).Str("function", fn).Str("tid", tid).Msg(msg)
}
func (l zeroLevelLogger) ErrorPretty(fn, tid string, msg string) {
	level := zerolog.Disabled
	if l.verbose {
		level = zerolog.ErrorLevel
	}
	l.PrettyErrLog.WithLevel(level).Str("function", fn).Str("tid", tid).Msg(msg)
}

func (l zeroLevelLogger) Print(level LogLevel, fn, tid string, msg string) {
	l.StdLog.Log().Str("level", getZeroLevel(level).String()).Str("function", fn).Str("tid", tid).Msg(msg)
}

func getZeroLevel(level LogLevel) zerolog.Level {
	switch level {
	case Info:
		return zerolog.InfoLevel
	case Warn:
		return zerolog.WarnLevel
	case Debug:
		return zerolog.DebugLevel
	case Error:
		return zerolog.ErrorLevel
	case Fatal:
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
