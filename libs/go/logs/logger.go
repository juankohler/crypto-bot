package logs

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type Format string

const (
	FormatJson   Format = "JSON"
	FormatText   Format = "TEXT"
	FormatPretty Format = "PRETTY"
)

type Level int

const (
	LevelDebug Level = Level(slog.LevelDebug)
	LevelError Level = Level(slog.LevelError)
	LevelInfo  Level = Level(slog.LevelInfo)
	LevelWarn  Level = Level(slog.LevelWarn)
)

type Config struct {
	LogLevel Level
	Format   Format
}

type logger struct {
	slog *slog.Logger
}

var loggerImp *logger = newLogger(&Config{})

func InitLogger(config *Config) {
	loggerImp = newLogger(config)
}

func newLogger(config *Config) *logger {
	logLevel := LevelDebug
	if config.LogLevel != 0 {
		logLevel = config.LogLevel
	}

	handlerOptions := &slog.HandlerOptions{
		Level: slog.Level(logLevel),
	}

	format := FormatJson
	if config.Format != "" {
		format = config.Format
	}

	var handler slog.Handler
	switch format {
	case FormatJson:
		handler = slog.NewJSONHandler(os.Stdout, handlerOptions)
	case FormatText:
		handler = slog.NewTextHandler(os.Stdout, handlerOptions)
	case FormatPretty:
		handler = newPrettyHandler(os.Stdout, PrettyHandlerOptions{
			SlogOpts: *handlerOptions,
		})
	default:
		handler = slog.NewJSONHandler(os.Stdout, handlerOptions)
	}

	return &logger{
		slog: slog.New(handler),
	}
}

func Info(ctx context.Context, msg string, attributes ...slog.Attr) {
	if loggerImp == nil {
		logNilLogger(msg, attributes...)
		return
	}

	collectorAttrs := getCollectorAttributes(ctx)
	allAttrs := append(collectorAttrs, attributes...)

	loggerImp.slog.LogAttrs(ctx, slog.LevelInfo, msg, allAttrs...)
}

func Warn(ctx context.Context, msg string, attributes ...slog.Attr) {
	if loggerImp == nil {
		logNilLogger(msg, attributes...)
		return
	}

	collectorAttrs := getCollectorAttributes(ctx)
	allAttrs := append(collectorAttrs, attributes...)

	loggerImp.slog.LogAttrs(ctx, slog.LevelWarn, msg, allAttrs...)
}

func Error(ctx context.Context, msg string, attributes ...slog.Attr) {
	if loggerImp == nil {
		logNilLogger(msg, attributes...)
		return
	}

	collectorAttrs := getCollectorAttributes(ctx)
	allAttrs := append(collectorAttrs, attributes...)

	loggerImp.slog.LogAttrs(ctx, slog.LevelError, msg, allAttrs...)
}

func Debug(ctx context.Context, msg string, attributes ...slog.Attr) {
	if loggerImp == nil {
		logNilLogger(msg, attributes...)
		return
	}

	collectorAttrs := getCollectorAttributes(ctx)
	allAttrs := append(collectorAttrs, attributes...)

	loggerImp.slog.LogAttrs(ctx, slog.LevelDebug, msg, allAttrs...)
}

func logNilLogger(msg string, attributes ...slog.Attr) {
	fmt.Println("WARNING: logger implementation is nil")
	fmt.Printf("message: %s \n", msg)
	for _, attr := range attributes {
		fmt.Printf("attribute: %s = %v\n", attr.Key, attr.Value)
	}
}
