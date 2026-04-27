package config

import (
	"github.com/RabotaemActivno/pulse/pkg/httpserver"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	Env string `envconfig:"ENV" defalult:"dev"`
}

//type HTTPConfig struct {
//	Port            string        `envconfig:"PORT" default:"8080"`
//	ReadTimeout     time.Duration `envconfig:"READ_TIMEOUT" default:"5s"`
//	WriteTimeout    time.Duration `envconfig:"WRITE_TIMEOUT" default:"10s"`
//	IdleTimeout     time.Duration `envconfig:"IDLE_TIMEOUT" default:"60s"`
//	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
//}

type Config struct {
	App  AppConfig
	HTTP httpserver.Config
}

func MustLoad() Config {
	var cfg Config
	if err := godotenv.Load(".env"); err != nil {
		panic("load config: " + err.Error())
	}
	if err := envconfig.Process("", &cfg); err != nil {
		panic("config: " + err.Error())
	}
	return cfg
}
