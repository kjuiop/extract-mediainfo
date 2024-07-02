package logger

import (
	"context"
	"extract-mediainfo/config"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
	"strings"
)

func SlogInit(cfg config.Logger) error {
	logLevel, err := slogLevelParser(cfg.Level)
	if err != nil {
		return err
	}

	var logWriter io.Writer
	if cfg.PrintStdOut {
		logWriter = io.MultiWriter(os.Stdout)
	} else {
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.Path,
			MaxSize:    100, // megabytes
			MaxBackups: 10,
			MaxAge:     28,    //days
			Compress:   false, // disabled by default
		}
		logWriter = fileWriter
	}

	handler := slog.NewJSONHandler(logWriter, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				s := a.Value.Any().(*slog.Source)
				s.File = getSimplePath(s.File)
			}
			return a
		},
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return nil
}

func slogLevelParser(lvStr string) (slog.Level, error) {
	dict := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
	result, ok := dict[strings.ToLower(lvStr)]
	if !ok {
		return result, fmt.Errorf("%s is not valid log level", lvStr)
	}
	return result, nil
}

func getSimplePath(path string) string {
	if len(path) <= 10 {
		return path
	}
	idx := strings.LastIndex(path, "/")
	if idx < 0 {
		return path
	}
	jdx := strings.LastIndex(path[:idx], "/")
	if jdx < 0 {
		return path
	}
	return path[jdx:]
}

func DefaultLogger(ctx context.Context) *slog.Logger {
	logger := slog.With("hostname", ctx.Value("hostname"))
	return logger
}
