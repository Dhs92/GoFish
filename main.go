package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/Dhs92/GoFish/config"
	"github.com/fsnotify/fsnotify"
)

func main() {
	// Setup logging
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	// Init config handler from Viper
	configHandler, err := config.ReadConfig()

	if err != nil {
		slog.Error("Error reading config", "error", err)
		os.Exit(1)
	}

	// Set log level from config
	logLevel, err := config.ParseLogLevel(configHandler.GetString("logLevel"))
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	if err != nil {
		slog.Error("Error parsing log level", "error", err)
	}

	slog.Info("Setting up logging", "logLevel", logLevel.String())

	slog.Debug("ListenAddr", "listenAddr", configHandler.GetString("listenAddr"))
	slog.Debug("ListenPort", "listenPort", configHandler.GetInt("listenPort"))
	slog.Debug("LogLevel:", "logLevel", configHandler.GetString("logLevel"))
	slog.Debug("Database.Host:", "database.host", configHandler.GetString("database.host"))
	slog.Debug("Database.Port:", "database.port", configHandler.GetInt("database.port"))
	slog.Debug("Database.User:", "database.user", configHandler.GetString("database.user"))
	slog.Debug("Database.Name:", "database.name", configHandler.GetString("database.name"))
	slog.Debug("Viper config: ", "config", configHandler.AllSettings())

	configHandler.OnConfigChange(func(e fsnotify.Event) {
		slog.Info("Config file changed:", "file", e.Name)
		logLevel, err := config.ParseLogLevel(configHandler.GetString("logLevel"))
		if err != nil {
			slog.Error("Error parsing log level", "error", err)
		}
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))
	})
	configHandler.WatchConfig()
	for {
		slog.Info("Sleeping for 5 seconds")
		slog.Debug("Debug message")
		slog.Warn("Warning message")
		slog.Error("Error message")
		slog.Info("Info message")
		time.Sleep(5 * time.Second)
	}
}
