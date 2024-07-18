package logger

import (
	"context"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

type MultiWriterHandler struct {
	stdoutHandler slog.Handler
	fileHandler   slog.Handler
	level         slog.Level
}

func (h *MultiWriterHandler) GetFileHandler() slog.Handler {
	return h.fileHandler
}

func (h *MultiWriterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *MultiWriterHandler) Handle(ctx context.Context, rec slog.Record) error {
	if err := h.stdoutHandler.Handle(ctx, rec); err != nil {
		return err
	}
	if err := h.fileHandler.Handle(ctx, rec); err != nil {
		return err
	}
	return nil
}

func (h *MultiWriterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &MultiWriterHandler{
		stdoutHandler: h.stdoutHandler.WithAttrs(attrs),
		fileHandler:   h.fileHandler.WithAttrs(attrs),
	}
}

func (h *MultiWriterHandler) WithGroup(name string) slog.Handler {
	return &MultiWriterHandler{
		stdoutHandler: h.stdoutHandler.WithGroup(name),
		fileHandler:   h.fileHandler.WithGroup(name),
	}
}

func NewMultiWriterHandler(logFile, output string) *MultiWriterHandler {
	// Configure lumberjack logger
	log := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,   // megabytes
		MaxBackups: 3,    // number of backup files
		MaxAge:     28,   // days
		Compress:   true, // compress old files
	}
	if output == "json" {
		return &MultiWriterHandler{
			stdoutHandler: slog.NewJSONHandler(os.Stdout, nil),
			fileHandler:   slog.NewJSONHandler(log, nil),
		}
	} else {
		return &MultiWriterHandler{
			stdoutHandler: slog.NewTextHandler(os.Stdout, nil),
			fileHandler:   slog.NewTextHandler(log, nil),
		}
	}
}

func New(logFile string, verbose, silent bool, output string) *slog.Logger {
	handler := NewMultiWriterHandler(logFile, output)
	if handler == nil {
		if output == "json" {
			handler = &MultiWriterHandler{
				stdoutHandler: slog.NewJSONHandler(os.Stdout, nil),
				fileHandler:   slog.Default().Handler(),
			}
		} else {
			handler = &MultiWriterHandler{
				stdoutHandler: slog.NewTextHandler(os.Stdout, nil),
				fileHandler:   slog.Default().Handler(),
			}
		}
	}

	if verbose {
		handler.level = slog.LevelDebug
	} else {
		handler.level = slog.LevelInfo
	}

	if silent {
		return slog.New(handler.GetFileHandler())
	}
	return slog.New(handler)
}
