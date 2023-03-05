package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug       bool `env:"IS_DEBUG" env-default:"false"`
	IsDevelopment bool `env:"IS_DEVELOP" env-default:"false"`
	Listen        struct {
		Type   string `env:"LISTEN_TYPE" env-default:"port"`
		BindIP string `env:"BIND_IP" env-default:"0.0.0.0"`
		Port   string `env:"PORT" env-default:"8080"`
	}
	Server struct {
		ReadTimeout  int    `env:"READ_TIMEOUT" env-default:"10"`
		WriteTimeout int    `env:"WRITE_TIMEOUT" env-default:"10"`
		IdleTimeout  int    `env:"IDLE_TIMEOUT" env-default:"100"`
		ServerType   string `env:"SERVER_TYPE" env-default:"http"`
	}
	AppConfig struct {
		LogLevel string `env:"LOGLEVEL" env-default:"debug"`
	}
	Database struct {
		Host     string `env:"DB_HOST" env-default:"localhost"`
		DbName   string `env:"DB_NAME" env-default:"anti-bruteforce-service-database"`
		Port     string `env:"DB_PORT" env-default:"5433"`
		User     string `env:"DB_USER" env-default:"postgres"`
		Password string `env:"DB_PASSWORD" env-default:"12345678"`
		SslMode  string `env:"SSL_MODE" env-default:"disable"`
	}
	Bucket struct {
		IpLimit             int `env:"IP_LIMIT" env-default:"1000"`
		LoginLimit          int `env:"LOGIN_LIMIT" env-default:"10"`
		PasswordLimit       int `env:"PASSWORD_LIMIT" env-default:"100"`
		ResetBucketInterval int `env:"RESET_BUCKET_INTERVAL" env-default:"60"`
	}
}

func New() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
