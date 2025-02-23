package config

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel zerolog.Level  `json:"logLevel"`
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"name"`
}

func defaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		LogLevel: zerolog.InfoLevel,
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "",
			Password: "",
			DbName:   "gofish",
		},
	}
}

func initViper() *viper.Viper {
	defaultConfig := defaultConfig()

	v := viper.New()

	// Set environment variables
	v.SetEnvPrefix("FISH")
	v.AutomaticEnv()

	// Set config file
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.SetConfigType("toml")
	log.Debug().Msgf("Config file name: %s", v.ConfigFileUsed())

	flagSet := pflag.NewFlagSet("config", pflag.ExitOnError)

	// Set flags
	flagSet.SetNormalizeFunc(wordSepNormalizeFunc)
	flagSet.String("server.host", defaultConfig.Server.Host, "Address to listen on")
	flagSet.Int("server.port", defaultConfig.Server.Port, "Port to listen on")
	flagSet.String("logLevel", defaultConfig.LogLevel.String(), "Log level")
	flagSet.String("database.host", defaultConfig.Database.Host, "Database host")
	flagSet.Int("database.port", defaultConfig.Database.Port, "Database port")
	flagSet.String("database.user", defaultConfig.Database.User, "Database user")
	flagSet.String("database.password", defaultConfig.Database.Password, "Database password")
	flagSet.Parse(os.Args[1:])

	// Bind flags to Viper
	v.BindPFlags(flagSet)
	// Set default values
	v.SetDefault("logLevel", defaultConfig.LogLevel)
	v.SetDefault("server.host", defaultConfig.Server.Host)
	v.SetDefault("server.port", defaultConfig.Server.Port)
	v.SetDefault("database.host", defaultConfig.Database.Host)
	v.SetDefault("database.port", defaultConfig.Database.Port)
	v.SetDefault("database.user", defaultConfig.Database.User)
	v.SetDefault("database.password", defaultConfig.Database.Password)
	v.SetDefault("database.name", defaultConfig.Database.DbName)

	return v
}

func ReadConfig() (*viper.Viper, error) {
	v := initViper()
	err := v.ReadInConfig()

	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Warn().Msg("Config file not found; creating a new one")
		v.SafeWriteConfig()
		err = nil
	} else if err != nil {
		log.Error().Err(err).Msg("Error reading config file")
		return nil, err
	}

	return v, err
}

func ParseLogLevel(level string) (zerolog.Level, error) {
	logLevel, err := zerolog.ParseLevel(level)
	return logLevel, err
}

func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	from := []string{"-", "_"}
	to := "."
	for _, sep := range from {
		name = strings.Replace(name, sep, to, -1)
		name = strings.ToLower(name)
	}
	return pflag.NormalizedName(name)
}

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
}
