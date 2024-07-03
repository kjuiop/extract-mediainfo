package main

import (
	"context"
	"extract-mediainfo/config"
	"extract-mediainfo/logger"
	"extract-mediainfo/watcher"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var BUILD_TIME = "no flag of BUILD_TIME"
var GIT_HASH = "no flag of GIT_HASH"
var APP_VERSION = "no flag of APP_VERSION"

func main() {

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	ctx = contextSetValue(ctx)

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("fail to read environments: %v", err)
	}

	if err := logger.SlogInit(cfg.Logger); err != nil {
		log.Fatalf("fail to init slog err : %v\n", err)
	}

	slog.Debug("extract-mediainfo start", "git_hash", GIT_HASH, "build_time", BUILD_TIME, "app_version", APP_VERSION, "hostname", ctx.Value("hostname"))

	watcherClient := watcher.NewWatcherClient(cfg.Watcher)

	wg.Add(1)
	go watcherClient.CheckWatcherFileCnt(ctx, &wg)

	<-exitSignal()
	cancel()
	wg.Wait()
	slog.Debug("gracefully stop mp-monitoring-app")
}

func exitSignal() <-chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	return sig
}

func contextSetValue(ctx context.Context) context.Context {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("fail load hostname : %v", err)
	}

	ctx = context.WithValue(ctx, "hostname", hostname)
	return ctx
}
