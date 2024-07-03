package watcher

import (
	"context"
	"extract-mediainfo/config"
	"extract-mediainfo/logger"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Watcher struct {
	cfg config.Watcher
}

func NewWatcherClient(cfg config.Watcher) *Watcher {
	return &Watcher{
		cfg: cfg,
	}
}

func (w *Watcher) CheckWatcherFileCnt(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	tick := time.Tick(time.Duration(w.cfg.CheckIntervalSec) * time.Second)

	for {
		select {
		case <-ctx.Done():
			slog.Debug("quit CheckWatcherFileCnt function")
			return
		case <-tick:
			msg := fmt.Sprintf("[%s] ", ctx.Value("hostname"))
			sLog := logger.DefaultLogger(ctx)

			for _, target := range w.cfg.Targets {
				count, err := countFiles(target.Path)
				if err != nil {
					sLog.Error("counts files err occur", "error", err)
					continue
				}

				sLog = sLog.With(target.Path, count)
				msg += fmt.Sprintf("[%s %d] ", target.Path, count)
			}
			sLog.Debug(msg)
		}
	}
}

func countFiles(dir string) (int, error) {
	var count int
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		count++
		return nil
	})
	return count, err
}
