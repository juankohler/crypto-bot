package logs

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"time"

	"github.com/juankohler/crypto-bot/libs/go/errors"
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()

	switch r.Level {
	case slog.LevelDebug:
		level = colorize(magenta, level)
	case slog.LevelInfo:
		level = colorize(green, level)
	case slog.LevelWarn:
		level = colorize(yellow, level)
	case slog.LevelError:
		level = colorize(red, level)
	}

	var fieldStr string
	r.Attrs(func(a slog.Attr) bool {
		val := a.Value.Any()

		if e, ok := val.(*errors.Error); ok {
			err, _ := e.MarshalJSON()
			fieldStr += fmt.Sprintf("%s=%+v ", a.Key, string(err))
			return true
		}

		fieldStr += fmt.Sprintf("%s=%+v ", a.Key, a.Value.Any())
		return true

	})

	timeStr := colorize(darkGray, r.Time.Format(time.DateTime))
	msg := r.Message

	h.l.Println(timeStr, level, msg, colorize(white, fieldStr))

	return nil
}

func newPrettyHandler(
	out io.Writer,
	opts PrettyHandlerOptions,
) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewTextHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}

	return h
}

const (
	reset        = "\033[0m"
	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	white        = 37
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
)

func colorize(colorCode int, text string) string {
	return fmt.Sprintf("\033[%dm%s%s", colorCode, text, reset)
}
