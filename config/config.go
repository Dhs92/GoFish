package config

type Config struct {
	ListenAddr string         `json:"listen_addr"`
	ListenPort int            `json:"listen_port"`
	LogLevel   string         `json:"log_level"`
	Database   DatabaseConfig `json:"database"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func DefaultConfig() Config {
	return Config{
		ListenAddr: "0.0.0.0",
		ListenPort: 8080,
		LogLevel:   "info",
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "",
			Password: "",
			Database: "gofish",
		},
	}
}
