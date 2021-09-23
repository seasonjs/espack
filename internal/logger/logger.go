package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"time"
)

type LogLevel uint8

const (
	LogLevelSilent LogLevel = iota
	LogLevelVerbose
	LogLevelDebug
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

type logger struct {
	log zerolog.Logger
}

func NewLogger() *logger {
	return &logger{log: zerolog.New(os.Stdout).With().Timestamp().Logger()}
}

// 设置需要的日志输出模式
func (l *logger) SetLogMode() *logger {
	return l
}

func (l *logger) Info(msg string, data ...interface{}) {
	if l.log.GetLevel() <= zerolog.InfoLevel {
		l.log.Info().Time("time", time.Now()).Msg(fmt.Sprintf(msg, data...))
	}
}

func (l *logger) Warn(msg string, data ...interface{}) {
	if l.log.GetLevel() <= zerolog.WarnLevel {
		l.log.Warn().Time("time", time.Now()).Msg(fmt.Sprintf(msg, data...))
	}
}

func (l *logger) Error(msg string, data ...interface{}) {
	if l.log.GetLevel() <= zerolog.ErrorLevel {
		l.log.Error().Time("time", time.Now()).Msg(fmt.Sprintf(msg, data...))
	}
}

func (l *logger) Trace(err error) {
	if err != nil {
		l.log.Error().Time("time", time.Now()).Msg(err.Error())
		return
	}

	level := l.log.GetLevel()

	if level <= zerolog.ErrorLevel {
	}
}
