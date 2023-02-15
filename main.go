package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thoughtgears/heimdall/internal/config"
	"github.com/thoughtgears/heimdall/internal/gcp"
	"github.com/thoughtgears/heimdall/internal/router"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.LevelFieldName = "severity"
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("config.New()")
	}

	client, err := gcp.New(cfg.Project)
	if err != nil {
		log.Fatal().Err(err).Msg("gcp.New(projectId)")
	}

	r := router.New(client, cfg)

	log.Info().Int("port", cfg.Port).Msg("launching server")

	if err := r.Run(fmt.Sprintf(":%v", cfg.Port)); err != nil {
		log.Fatal().Err(err).Msg("r.Run()")
	}
}
