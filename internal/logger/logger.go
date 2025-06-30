package logger

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

// Init инициализирует глобальный логгер
func Init() {
	// Создаем JSON handler для структурированного логирования
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		AddSource: true,
	})
	
	logger = slog.New(handler)
	slog.SetDefault(logger)
}

// Get возвращает экземпляр логгера
func Get() *slog.Logger {
	if logger == nil {
		Init()
	}
	return logger
}

// Info логирует информационное сообщение
func Info(msg string, args ...any) {
	Get().Info(msg, args...)
}

// Error логирует ошибку
func Error(msg string, args ...any) {
	Get().Error(msg, args...)
}

// Debug логирует отладочное сообщение
func Debug(msg string, args ...any) {
	Get().Debug(msg, args...)
}

// Warn логирует предупреждение
func Warn(msg string, args ...any) {
	Get().Warn(msg, args...)
}

// Fatal логирует критическую ошибку и завершает программу
func Fatal(msg string, args ...any) {
	Get().Error(msg, args...)
	os.Exit(1)
}

// With добавляет атрибуты к логгеру
func With(args ...any) *slog.Logger {
	return Get().With(args...)
}
