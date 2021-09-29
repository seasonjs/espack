package logger

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"seasonjs/espack/internal/utils"
	"sync"
)

var (
	once sync.Once
	log  zerolog.Logger
)

func init() {
	once.Do(func() {
		newLogger()
		if utils.Env.Dev() {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}

	})
}

//time.FC3339     = "2006-01-02T15:04:05Z07:00"
func newLogger() {
	zerolog.ErrorStackFieldName = "call stack"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	//time.FC3339     = "2006-01-02T15:04:05Z07:00"
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04"}
	output.FormatLevel = func(i interface{}) string {
		return fmt.Sprintf("[espack %s]", i)
	}

	//output.FormatMessage = func(i interface{}) string {
	//	return fmt.Sprintf("***%s****", i)
	//	//return ""
	//}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
		//return ""
	}
	//output.FormatFieldValue= func(i interface{}) string {
	//	return fmt.Sprintf("%s:", i)
	//}
	output.FormatErrFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
		//return ""
	}
	output.FormatErrFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s\n", i)
		//return ""
	}
	output.PartsExclude = append(output.PartsExclude, "time")
	log = zerolog.New(output).With().Timestamp().Logger()
}

func Info(format string, data ...interface{}) {
	log.Info().Msgf(format, data...)
}

func Warn(format string, data ...interface{}) {
	log.Warn().Msgf(format, data...)
}

func Error(format string, data ...interface{}) {

	log.Error().Msgf(format, data...)
}

func Fail(err error, msg string) {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	e, ok := err.(stackTracer)
	if !ok {
		return
	}
	if utils.Env.Dev() {
		log.Fatal().Err(err).Msg(msg)
		for _, frame := range e.StackTrace() {
			fmt.Printf("%+s:%d\r\n", frame, frame)
		}
	} else {
		log.Fatal().Err(err).Msg(msg)
	}
}

func Trace(err error, msg string) {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	e, ok := err.(stackTracer)
	if !ok {
		return
	}
	//It's mean when env=dev just print track
	if utils.Env.Dev() {
		log.Error().Msg(msg)
		for _, frame := range e.StackTrace() {
			fmt.Printf("%+s:%d\r\n", frame, frame)
		}
	} else {
		log.Error().Stack().Err(err).Msg(msg)
	}
}
