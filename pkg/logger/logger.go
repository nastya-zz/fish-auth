package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var logger *slog.Logger

// Цветовые коды для ANSI
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorGreen  = "\033[32m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[90m"
	ColorBold   = "\033[1m"
)

// PrettyHandler - кастомный хендлер для красивого вывода в режиме разработки
type PrettyHandler struct {
	writer io.Writer
	level  slog.Level
}

// NewPrettyHandler создает новый экземпляр PrettyHandler
func NewPrettyHandler(w io.Writer, level slog.Level) *PrettyHandler {
	return &PrettyHandler{
		writer: w,
		level:  level,
	}
}

// Enabled проверяет, нужно ли логировать сообщение данного уровня
func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle обрабатывает лог-запись
func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	// Форматируем время
	timeStr := r.Time.Format("15:04:05.000")
	
	// Определяем цвет и префикс для уровня логирования
	var levelColor, levelStr string
	switch r.Level {
	case slog.LevelDebug:
		levelColor = ColorGray
		levelStr = "DEBUG"
	case slog.LevelInfo:
		levelColor = ColorBlue
		levelStr = "INFO "
	case slog.LevelWarn:
		levelColor = ColorYellow
		levelStr = "WARN "
	case slog.LevelError:
		levelColor = ColorRed
		levelStr = "ERROR"
	default:
		levelColor = ColorReset
		levelStr = "UNKN "
	}
	
	// Строим основную часть сообщения
	var output strings.Builder
	
	// Время (серым цветом)
	output.WriteString(fmt.Sprintf("%s%s%s ", ColorGray, timeStr, ColorReset))
	
	// Уровень логирования (цветной)
	output.WriteString(fmt.Sprintf("%s%s%s%s ", ColorBold, levelColor, levelStr, ColorReset))
	
	// Сообщение (жирным шрифтом для важных уровней)
	if r.Level >= slog.LevelWarn {
		output.WriteString(fmt.Sprintf("%s%s%s", ColorBold, r.Message, ColorReset))
	} else {
		output.WriteString(r.Message)
	}
	
	// Обрабатываем атрибуты
	var attrs []string
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, formatAttr(a))
		return true
	})
	
	if len(attrs) > 0 {
		output.WriteString(fmt.Sprintf(" %s[%s]%s", ColorCyan, strings.Join(attrs, ", "), ColorReset))
	}
	
	output.WriteString("\n")
	
	_, err := h.writer.Write([]byte(output.String()))
	return err
}

// WithAttrs добавляет атрибуты к хендлеру
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Для простоты возвращаем тот же хендлер
	// В более сложной реализации можно сохранять атрибуты
	return h
}

// WithGroup добавляет группу к хендлеру
func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	// Для простоты возвращаем тот же хендлер
	return h
}

// formatAttr форматирует атрибут для вывода
func formatAttr(attr slog.Attr) string {
	switch attr.Value.Kind() {
	case slog.KindString:
		return fmt.Sprintf("%s=%q", attr.Key, attr.Value.String())
	case slog.KindTime:
		return fmt.Sprintf("%s=%s", attr.Key, attr.Value.Time().Format(time.RFC3339))
	default:
		return fmt.Sprintf("%s=%v", attr.Key, attr.Value.Any())
	}
}

// isDevelopment проверяет, запущено ли приложение в режиме разработки
func isDevelopment() bool {
	env := os.Getenv("ENV")
	if env == "" {
		env = os.Getenv("ENVIRONMENT")
	}
	if env == "" {
		env = os.Getenv("GO_ENV")
	}
	return env == "" || env == "dev" || env == "development" || env == "local"
}

// Init инициализирует глобальный логгер
func Init() {
	var handler slog.Handler
	
	if isDevelopment() {
		// В режиме разработки используем красивый вывод
		handler = NewPrettyHandler(os.Stdout, slog.LevelDebug)
	} else {
		// В продакшене используем JSON handler
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
		})
	}
	
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

// InfoWithCaller логирует информационное сообщение с информацией о вызывающем коде
func InfoWithCaller(msg string, args ...any) {
	if isDevelopment() {
		caller := getCaller(2)
		args = append(args, "caller", caller)
	}
	Get().Info(msg, args...)
}

// ErrorWithCaller логирует ошибку с информацией о вызывающем коде
func ErrorWithCaller(msg string, args ...any) {
	if isDevelopment() {
		caller := getCaller(2)
		args = append(args, "caller", caller)
	}
	Get().Error(msg, args...)
}

// DebugWithCaller логирует отладочное сообщение с информацией о вызывающем коде
func DebugWithCaller(msg string, args ...any) {
	if isDevelopment() {
		caller := getCaller(2)
		args = append(args, "caller", caller)
	}
	Get().Debug(msg, args...)
}

// WarnWithCaller логирует предупреждение с информацией о вызывающем коде
func WarnWithCaller(msg string, args ...any) {
	if isDevelopment() {
		caller := getCaller(2)
		args = append(args, "caller", caller)
	}
	Get().Warn(msg, args...)
}

// getCaller возвращает информацию о вызывающем коде
func getCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// SetLevel устанавливает уровень логирования
func SetLevel(level slog.Level) {
	if logger != nil {
		// Пересоздаем логгер с новым уровнем
		var handler slog.Handler
		if isDevelopment() {
			handler = NewPrettyHandler(os.Stdout, level)
		} else {
			handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level:     level,
				AddSource: true,
			})
		}
		logger = slog.New(handler)
		slog.SetDefault(logger)
	}
}
