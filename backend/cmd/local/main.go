package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/saddiqs1/gymshark-saddiqs1/config"
	"github.com/saddiqs1/gymshark-saddiqs1/internal/api"
	"github.com/saddiqs1/gymshark-saddiqs1/pkg/logger"
)

func main() {
	cfg, err := config.NewLocalConfig()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	logger := logger.New(cfg.Environment.LogLevel, cfg.Environment.AppEnv, nil)
	server := &http.Server{
		Addr:    ":8080",
		Handler: api.NewRouter(),
	}

	logger.Info().Msgf("server running on %s", server.Addr)

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal().Err(err).Msg("server stopped")
	}
}
