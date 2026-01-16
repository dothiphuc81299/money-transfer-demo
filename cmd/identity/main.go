package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"money-transfer-demo/cmd/identity/app"
	"money-transfer-demo/pkg/infra/log"

	"go.uber.org/zap"
)

func main() {
	zlog, err := log.New("identity-main")
	if err != nil {
		panic("failed to init logger: " + err.Error())
	}
	
	defer zlog.Sync()

	zlog.Info("Starting Identity Service...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		zlog.Warn("Received shutdown signal", zap.String("signal", sig.String()))
		cancel()
	}()

	server, err := app.NewServer()
	if err != nil {
		zlog.Fatal("Failed to initialize server", zap.Error(err))
	}

	<-ctx.Done()

	zlog.Info("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		zlog.Fatal("Failed to shutdown server", zap.Error(err))
	}

	zlog.Info("Server exited gracefully")
}
