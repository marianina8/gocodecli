package logger

import (
	"context"
	"log/slog"
	"os"
)

type MultiWriterHandler struct {
	stdoutHandler slog.Handler
	fileHandler   slog.Handler
}

func (h *MultiWriterHandler) GetFileHandler() slog.Handler {
	return h.fileHandler
}

func (h *MultiWriterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *MultiWriterHandler) Handle(ctx context.Context, rec slog.Record) error {
	if err := h.stdoutHandler.Handle(ctx, rec); err!= nil {
		return err
	}
	if err := h.fileHandler.Handle(ctx, rec); err!= nil {
		return err
	}
	return nil
}

func (h *MultiWriterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &MultiWriterHandler{
		stdoutHandler: h.stdoutHandler.WithAttrs(attrs),
		fileHandler: h.fileHandler.WithAttrs(attrs),
	}
}

func (h *MultiWriterHandler) WithGroup(name string) slog.Handler {
	return &MultiWriterHandler{
		stdoutHandler: h.stdoutHandler.WithGroup(name),
		fileHandler: h.fileHandler.WithGroup(name),
	}
}

func NewMultiWriterHandler(logFile string) *MultiWriterHandler {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Filed to open log file", "error", err, "path", logFile)
		return nil
	}
	return &MultiWriterHandler{
		stdoutHandler: slog.NewTextHandler(os.Stdout, nil),
		fileHandler: slog.NewTextHandler(file, nil),
	}
}

func New(logFile string, verbose, silent bool) *slog.Logger {
	handler:= NewMultiWriterHandler(logFile)
	if handler == nil {
		handler = &MultiWriterHandler {
			stdoutHandler: slog.NewTextHandler(os.Stdout, nil),
			fileHandler: slog.Default().Handler(),
		}
	}

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	if silent {
		return slog.New(handler.GetFileHandler())
	}
	return slog.New(handler)
}

