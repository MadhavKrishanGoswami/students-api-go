package main

import (
	"log"
	"net/http"

	"github.com/MadhavKrishanGoswami/students-api/internal/config"
)

func main() {
	// 1. Load config
	cfg := config.MustLoad()

	// 2. Setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to students-api"))
	})

	// 3. Setup Server
	log.Printf("Starting server on %s", cfg.Addr)
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	// 4. Start the server
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
