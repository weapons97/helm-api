package utils

import (
	"github.com/rs/zerolog"
	"os"
)

var (
	loge = zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
	logi = zerolog.New(os.Stderr).With().Stack().Timestamp().Logger()
)

func Debug() *zerolog.Event {
	return loge.Debug()
}
func Info() *zerolog.Event {
	return logi.Info()
}
func Warn() *zerolog.Event {
	return loge.Warn()
}
func Error() *zerolog.Event {
	return loge.Error()
}
func Panic() *zerolog.Event {
	return loge.Panic()
}
