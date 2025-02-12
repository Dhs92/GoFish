package config

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel slog.Level     `json:"logLevel"`
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
		LogLevel: slog.LevelInfo,
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "",
			Password: "",
			DbName:   "gofish",
		},
	}
}

func InitViper() *viper.Viper {
	defaultConfig := defaultConfig()

	v := viper.New()

	// Set environment variables
	v.SetEnvPrefix("FISH")
	v.AutomaticEnv()

	// Set config file
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.SetConfigType("toml")
	slog.Debug("Config file name", "value", v.ConfigFileUsed())

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
	v := InitViper()
	err := v.ReadInConfig()

	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		slog.Warn("Config file not found; creating a new one")
		v.SafeWriteConfig()
		err = nil
	} else if err != nil {
		slog.Error("Error reading config file", "error", err)
		return nil, err
	}

	return v, err
}

func ParseLogLevel(level string) (slog.Level, error) {
	var logLevel slog.Level
	err := logLevel.UnmarshalText([]byte(level))
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
