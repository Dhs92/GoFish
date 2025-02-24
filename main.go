package main

import (
	"os"
	"time"

	"github.com/Dhs92/GoFish/config"
	"github.com/Dhs92/GoFish/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Setup logging
	logger.InitLogger()

	// Init config handler from Viper
	configHandler, err := config.ReadConfig()

	if err != nil {
		log.Error().Err(err).Msg("Error reading config")
		os.Exit(1)
	}

	// Set log level from config
	logLevel, err := config.ParseLogLevel(configHandler.GetString(config.LOGLEVEL))
	if err != nil {
		log.Error().Err(err).Msg("Error parsing log level")
	} else {
		zerolog.SetGlobalLevel(logLevel)
	}

	log.Info().Str("logLevel", logLevel.String()).Msg("Setting up logging")

	log.Debug().Str("server.host", configHandler.GetString(config.SERVERHOST)).Msg("")
	log.Debug().Int("server.port", configHandler.GetInt(config.SERVERPORT)).Msg("")
	log.Debug().Str("logLevel", configHandler.GetString(config.LOGLEVEL)).Msg("")
	log.Debug().Str("database.host", configHandler.GetString(config.DATABASEHOST)).Msg("")
	log.Debug().Int("database.port", configHandler.GetInt(config.DATABASEPORT)).Msg("")
	log.Debug().Str("database.user", configHandler.GetString(config.DATABASEUSER)).Msg("")
	log.Debug().Str("database.name", configHandler.GetString(config.DATABASENAME)).Msg("")
	log.Debug().Interface("config", configHandler.AllSettings()).Msg("")

	configHandler.OnConfigChange(func(e fsnotify.Event) {
		log.Info().Str("file", e.Name).Msg("Config file changed")
		logLevel, err := config.ParseLogLevel(configHandler.GetString(config.LOGLEVEL))
		if err != nil {
			log.Error().Err(err).Msg("Error parsing log level")
		} else {
			log.Debug().Str(config.LOGLEVEL, logLevel.String()).Msg("Setting log level")
			zerolog.SetGlobalLevel(logLevel)
		}
	})
	configHandler.WatchConfig()

	for {
		log.Info().Msg("Sleeping for 5 seconds")
		log.Debug().Msg("Debug message")
		log.Warn().Msg("Warning message")
		log.Error().Msg("Error message")
		log.Info().Msg("Info message")
		time.Sleep(5 * time.Second)
	}
}
