package logger

import (
	"log/slog"
	"os"
)

// Определение структуры для логгера Slog
type Logger struct {
	Json *slog.Logger
}

// Определение внешней переменной логгера для всего проекта
var Slog *Logger

// Метод инициализации логгера
func Init() {
	Slog = &Logger{Json: slog.New(slog.NewJSONHandler(os.Stdout, nil))}
}
