package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	JwtAuth    `yaml:"jwt_auth"`
	HTTPServer `yaml:"http_server"`
	DataBase   `yaml:"database"`
}

type JwtAuth struct {
	PrivateKeyPath string        `yaml:"private_key_path" env:"PRIVATE_KEY_PATH"`
	PublicKeyPath  string        `yaml:"public_key_path" env:"PUBLIC_KEY_PATH"`
	JwtTTL         time.Duration `yaml:"jwt_ttl" env:"JWT_TTL"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DataBase struct {
	Type     string `yaml:"type" env-default:"postgres" ENV:"DATABASE_TYPE"`
	Name     string `yaml:"name" env:"DATABASE_NAME"`
	Address  string `yaml:"address" env:"DATABASE_ADDRESS"`
	Username string `yaml:"username" env:"DATABASE_USERNAME"`
	Password string `yaml:"password" env:"DATABASE_PASSWORD"`
}

func NewConfig(env string) *Config {
	if env == "local" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// get env variables
	configPath := os.Getenv("CONFIG_PATH")

	// check if file exists
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &config
}
