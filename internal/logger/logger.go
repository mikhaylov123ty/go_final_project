package logger

import (
	"log"
	"log/slog"
	"os"
)

type Logger struct {
	JsonInfo  *log.Logger
	JsonWarn  *log.Logger
	JsonError *log.Logger
}

var Slog *Logger

func Init() {
	Slog = &Logger{
		JsonInfo:  slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), 0),
		JsonWarn:  slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), 4),
		JsonError: slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), 8),
	}
}
