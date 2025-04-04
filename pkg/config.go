package pkg

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config keys
const (
	LogLevel         string = "logLevel"
	LogFormat        string = "logFormat" // "json" or "pretty"
	ServerHost       string = "server.host"
	ServerPort       string = "server.port"
	DatabaseDriver   string = "database.driver" // e.g., "postgres", "sqlite"
	DatabaseHost     string = "database.host"
	DatabasePort     string = "database.port"
	DatabaseUser     string = "database.user"
	DatabasePassword string = "database.password"
	DatabaseName     string = "database.name"
)

type Config struct {
	LogLevel  string         `toml:"logLevel"`
	LogFormat string         `toml:"logFormat"` // "json" or "pretty"
	Server    ServerConfig   `toml:"server"`
	Database  DatabaseConfig `toml:"database"`
}

type ServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type DatabaseConfig struct {
	Driver   string `toml:"driver"` // e.g., "postgres", "sqlite"
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Name     string `toml:"name"`
}

func defaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		LogLevel:  zerolog.InfoLevel.String(),
		LogFormat: "json", // Default to JSON format for structured logging
		Database: DatabaseConfig{
			Driver:   "sqlite", // Change to your preferred driver
			Host:     "localhost",
			Port:     5432,
			User:     "gofish",
			Password: "gofish",
			Name:     "gofish",
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

	flagSet := pflag.NewFlagSet("config", pflag.ExitOnError)

	// Set flags
	flagSet.SetNormalizeFunc(wordSepNormalizeFunc)
	flagSet.String(ServerHost, defaultConfig.Server.Host, "Address to listen on")
	flagSet.Int(ServerPort, defaultConfig.Server.Port, "Port to listen on")
	flagSet.String(LogLevel, defaultConfig.LogLevel, "Log level")
	flagSet.String(LogFormat, defaultConfig.LogFormat, "Log format (toml or pretty)")
	// Database flags
	flagSet.String(DatabaseDriver, defaultConfig.Database.Driver, "Database driver (e.g., postgres, sqlite)")
	flagSet.String(DatabaseHost, defaultConfig.Database.Host, "Database host")
	flagSet.Int(DatabasePort, defaultConfig.Database.Port, "Database port")
	flagSet.String(DatabaseUser, defaultConfig.Database.User, "Database user")
	flagSet.String(DatabasePassword, defaultConfig.Database.Password, "Database password")
	flagSet.Parse(os.Args[1:])

	// Bind flags to Viper
	v.BindPFlags(flagSet)
	// Set default values
	v.SetDefault(LogLevel, defaultConfig.LogLevel)
	v.SetDefault(LogFormat, defaultConfig.LogFormat)
	// Set default values for server and database
	v.SetDefault(ServerHost, defaultConfig.Server.Host)
	v.SetDefault(ServerPort, defaultConfig.Server.Port)
	// Database defaults
	v.SetDefault(DatabaseDriver, defaultConfig.Database.Driver)
	v.SetDefault(DatabaseHost, defaultConfig.Database.Host)
	v.SetDefault(DatabasePort, defaultConfig.Database.Port)
	v.SetDefault(DatabaseUser, defaultConfig.Database.User)
	v.SetDefault(DatabasePassword, defaultConfig.Database.Password)
	v.SetDefault(DatabaseName, defaultConfig.Database.Name)

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

func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	from := []string{"-", "_"}
	to := "."
	for _, sep := range from {
		name = strings.Replace(name, sep, to, -1)
		name = strings.ToLower(name)
	}
	return pflag.NormalizedName(name)
}
