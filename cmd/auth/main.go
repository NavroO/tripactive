package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NavroO/tripactive/internal/auth"
	"github.com/NavroO/tripactive/internal/shared"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func main() {
	shared.SetupLogger()
	log.Info().Msg("logger initialized")
	cfg := shared.LoadConfig()
	if cfg.Port == "" {
		log.Fatal().Msg("PORT is not set in .env")
	}
	db, err := shared.ConnectDB()
	if err != nil {
		log.Fatal().Msgf("cannot connect to db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close database")
		}
	}()

	r := chi.NewRouter()
	svc := auth.NewService()
	h := auth.NewHandler(svc)
	r.Mount("/auth", h.Routes())

	srv := &http.Server{
		Addr:         "0.0.0.0:" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Msgf("starting HTTP server on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Msgf("server forced to shutdown: %v", err)
	}

	log.Info().Msg("server exited gracefully")
}
