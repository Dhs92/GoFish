package main

import (
	"fmt"
	"github.com/Dhs92/GoFish/internal/models"
	"github.com/Dhs92/GoFish/internal/repository"
	"github.com/Dhs92/GoFish/internal/service"
	"github.com/Dhs92/GoFish/pkg"
	"github.com/fsnotify/fsnotify"
	"github.com/glebarez/sqlite" // Importing the SQLite driver for GORM
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	// Importing the logger package from GORM
)

func main() {
	// Init config handler from Viper
	configHandler, err := pkg.ReadConfig()

	if err != nil {
		log.Error().Err(err).Msg("Error reading config")
		os.Exit(1)
	}

	var configFile pkg.Config

	err = configHandler.Unmarshal(&configFile)

	if err != nil {
		log.Error().Err(err).Msg("Error unmarshalling config")
		os.Exit(1)
	}

	// Setup logging
	pkg.InitLogger(configFile.LogFormat)

	logLevel, err := zerolog.ParseLevel(configFile.LogLevel)

	if err != nil {
		log.Error().Err(err).Msg("Error parsing log level")
		os.Exit(1)
	}
	zerolog.SetGlobalLevel(logLevel) // Set the global log level

	log.Info().Str("logLevel", configFile.LogLevel).Msg("Setting up logging")

	//TODO: Condense the database connection logic into a single function

	var conn gorm.Dialector
	if configFile.Database.Driver == "sqlite" {
		conn = sqlite.Open("gofish.db")
	} else if configFile.Database.Driver == "postgres" {
		conn = postgres.Open(fmt.Sprintf("postgres://%s:%s@%s:%v/%s", configFile.Database.User, configFile.Database.Password, configFile.Database.Host, configFile.Database.Port, configFile.Database.Name))
	}
	db, err := gorm.Open(conn, &gorm.Config{Logger: pkg.NewGormLogger()})

	// GORM will log the error if the connection fails, this will bail if the connection fails
	if err != nil {
		os.Exit(1)
	}

	// Automatically migrate the database schema for the models
	// TODO: Condense the migration logic into a single function
	err = db.AutoMigrate(&models.User{}, &models.Livestock{}, &models.Tank{}, &models.StockItem{}, &models.ScheduleItem{})
	if err != nil {
		log.Error().Err(err).Msg("Error auto migrate")
		os.Exit(1)
	}

	userService := service.NewUserService(repository.NewGormUserRepository(db), log.Logger)
	newUserID, _ := userService.CreateUser("John", "john@test.com", "password123", models.RoleUser)

	// TODO: Create tank service
	db.Save(models.NewTank(*newUserID, "Main Tank", 100, "gal"))
	db.Save(models.NewTank(*newUserID, "Main Tank", 100, "gal"))

	configHandler.OnConfigChange(func(e fsnotify.Event) {
		err = configHandler.ReadInConfig()
		if err != nil {
			log.Error().Err(err).Msg("Error reading config after change")
		}
		err = configHandler.Unmarshal(&configFile) // Reload the config file
		if err != nil {
			log.Error().Err(err).Msg("Error unmarshalling config after change")
			return
		}

		log.Info().Str("file", e.Name).Msg("Config file changed")
		logLevel, err := zerolog.ParseLevel(configFile.LogLevel)
		if err != nil {
			log.Error().Err(err).Msg("Error parsing log level")
		} else {
			log.Debug().Str("logLevel", pkg.LogLevel).Msg("Setting log level")
			zerolog.SetGlobalLevel(logLevel)
		}
	})
	configHandler.WatchConfig()
}
