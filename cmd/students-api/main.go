package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MadhavKrishanGoswami/students-api/internal/config"
	"github.com/MadhavKrishanGoswami/students-api/internal/http/handlers/student"
)

func main() {
	// 1. Load config
	cfg := config.MustLoad()

	// 2. Setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New())

	// 3. Setup Server
	slog.Info("server started", "address", cfg.Addr)
	log.Printf("Starting server on %s", cfg.Addr)
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	// 4. Start the server
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()
	<-done

	slog.Info("Shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server", "error", err)
	}
	slog.Info("Server gracefully stopped")
}
