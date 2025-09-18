package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wallet-app/internal/config"
	"wallet-app/internal/db"
	"wallet-app/internal/handlers"
	"wallet-app/internal/repo"
	"wallet-app/internal/server"
	"wallet-app/internal/service"
)

func main() {
	cfg := config.Load()

	pool, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	repo := repo.NewWalletRepo(pool)
	svc := service.NewWalletService(repo)
	handler := handlers.NewWalletHandler(svc)

	srv := server.New(handler, cfg.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != server.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}
