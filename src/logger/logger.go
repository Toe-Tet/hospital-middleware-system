package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func Init(level string) {
	var lvl zerolog.Level
	switch level {
	case "production", "prod":
		lvl = zerolog.InfoLevel
	case "test":
		lvl = zerolog.ErrorLevel
	default:
		lvl = zerolog.DebugLevel
	}

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	Log = zerolog.New(output).
		Level(lvl).
		With().
		Timestamp().
		Caller().
		Logger()
}
