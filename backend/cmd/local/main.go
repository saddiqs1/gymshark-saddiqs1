package main

import (
	"log"

	"github.com/saddiqs1/gymshark-saddiqs1/config"
	"github.com/saddiqs1/gymshark-saddiqs1/pkg/logger"
)

func main() {
	cfg, err := config.NewLocalConfig()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	logger := logger.New(cfg.Environment.LogLevel, cfg.Environment.AppEnv, nil)
	logger.Info().Msgf("hello world - local")

	// TODO - spin up local api using net/http
	// local.Run(context.Background(), logger, cfg)
}
