package main

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	// Initialize and start the SmartHome Monitor

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	debugEnv := os.Getenv("DEBUG")
	debug := flag.Bool("debug", debugEnv == "true", "sets log level to debug")

	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Info().Msg("Starting SmartHome Monitor")
	log.Info().Msgf("Debug mode enabled: %t", *debug)

	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config.yaml"
	}

	monitor, err := NewMonitor("config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing monitor")

	}

	monitor.Start()
}
