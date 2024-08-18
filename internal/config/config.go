package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string `yaml:"env" env-default:"local"`
	HTTPServer   `yaml:"http_server"`
	DatabaseData `yaml:"database_data"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env:"HOST_PORT" env-default:":8082"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`

	//User     string `yaml:"user" env-required:"true"`
	//Password string `yaml:"password" env-required:"true"`
}

type DatabaseData struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST" env-default:"localhost"`
	Port     int64  `yaml:"port" env:"POSTGRES_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"POSTGRES_USER" env-default:"test"`
	Password string `yaml:"password" env-required:"true" env:"POSTGRES_PASSWORD" env-default:"test"`
	DBName   string `yaml:"dbname" env:"POSTGRES_DB" env-default:"postgres"`
}

func MustLoad() *Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	fmt.Println(cfg)
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
