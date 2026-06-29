package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Rowkash/go-gin-auth/internal/app"
	"github.com/Rowkash/go-gin-auth/internal/config"
	"github.com/caarlos0/env/v11"
)

func main() {

	var cfg config.Config

	if err := env.Parse(&cfg); err != nil {
		log.Printf("Invalid Environment Variables: %v", err)
		return
	}

	application := app.NewApp(cfg)

	go func() {
		if err := application.Run(":" + cfg.App.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Listen error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
