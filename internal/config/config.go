package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string `yaml:"env" env-default:"local"`
	HTTPServer   `yaml:"http_server"`
	DatabaseData `yaml:"database_data"`
	//EmailData    `yaml:"email_data"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8082"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`

	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type DatabaseData struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int64  `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"test"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"dbname" env-default:"postgres"`
}

//type EmailData struct {
//	EmailSenderName     string `yaml:"EMAIL_SENDER_NAME"`
//	EmailSenderAddress  string `yaml:"EMAIL_SENDER_ADDRESS"`
//	EmailSenderPassword string `yaml:"EMAIL_SENDER_PASSWORD"`
//}

func MustLoad() *Config {
	err := os.Setenv("CONFIG_PATH", "./config/local.yaml")
	if err != nil {
		log.Fatal("fail to set CONFIG_PATH")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesn't exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
