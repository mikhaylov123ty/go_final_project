package logger

import (
	"log"
	"log/slog"
	"os"
)

// Определение структуры для логгера Slog
type Logger struct {
	JsonInfo  *log.Logger
	JsonWarn  *log.Logger
	JsonError *log.Logger
}

// Определение внешней переменной логеера для всего проекта
var Slog *Logger

// Метод инициализации логгера
func Init() {
	Slog = &Logger{
		JsonInfo:  slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), 0),
		JsonWarn:  slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), 4),
		JsonError: slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), 8),
	}
}
