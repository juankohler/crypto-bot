package logs

import (
	"context"
	"fmt"
	"log/slog"
)

type loggerCollector struct {
	attributes []slog.Attr
}

type loggerCollectorKey struct{}

func NewAttr(key string, value interface{}) slog.Attr {
	return slog.Any(key, value)
}

func AddAttr(ctx context.Context, attributes ...slog.Attr) context.Context {
	collector := fromContext(ctx)
	if collector == nil {
		fmt.Println("Add log attr ignored because of nil collector")
		return ctx
	}

	collector.attributes = append(collector.attributes, attributes...)
	return toContext(ctx, collector)
}

func getCollectorAttributes(ctx context.Context) []slog.Attr {
	collector := fromContext(ctx)
	if collector == nil {
		loggerImp.slog.Warn("The collector logger is nil")
		return []slog.Attr{}
	}
	return collector.attributes
}

func ContextWithLogger(ctx context.Context) context.Context {
	return toContext(ctx, &loggerCollector{})
}

func fromContext(ctx context.Context) *loggerCollector {
	if collector, ok := ctx.Value(loggerCollectorKey{}).(*loggerCollector); ok {
		return collector
	}
	return nil
}

func toContext(ctx context.Context, rCollector *loggerCollector) context.Context {
	return context.WithValue(ctx, loggerCollectorKey{}, rCollector)
}
