package config

import (
	"github.com/RabotaemActivno/pulse/pkg/httpserver"
	"github.com/RabotaemActivno/pulse/pkg/logger"
	"github.com/RabotaemActivno/pulse/pkg/postgres"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	Env string `envconfig:"ENV" default:"dev"`
}

//type HTTPConfig struct {
//	Port            string        `envconfig:"PORT" default:"8080"`
//	ReadTimeout     time.Duration `envconfig:"READ_TIMEOUT" default:"5s"`
//	WriteTimeout    time.Duration `envconfig:"WRITE_TIMEOUT" default:"10s"`
//	IdleTimeout     time.Duration `envconfig:"IDLE_TIMEOUT" default:"60s"`
//	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
//}

type Config struct {
	App      AppConfig
	HTTP     httpserver.Config
	Logger   logger.Config
	Postgres postgres.Config
}

func MustLoad() Config {
	var cfg Config
	_ = godotenv.Load(".env")
	if err := envconfig.Process("", &cfg); err != nil {
		panic("config: " + err.Error())
	}
	return cfg
}
