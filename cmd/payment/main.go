package main

import (
	"context"
	"log"
	"money-transfer-demo/cmd/payment/app"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		sig := <-sigChan
		log.Printf("⚠️ Received signal: %v. Shutting down...\n", sig)
		cancel()
	}()

	server, err := app.NewServer()
	if err != nil {
		log.Fatalf("❌ Failed to initialize server: %v", err)
	}

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("❌ Failed to shutdown server: %v", err)
	}

	sqlDB, err := server.Postgresdb.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get SQL DB: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("❌ Failed to close database: %v", err)
	}

	log.Println("✅ Server shutdown gracefully")
}
