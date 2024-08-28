package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string        `yaml:"env" env-default:"local"`
	JwtSecret  string        `yaml:"jwt_secret"`
	JwtTTL     time.Duration `yaml:"jwt_ttl"`
	HTTPServer `yaml:"http_server"`
	DataBase   `yaml:"database"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DataBase struct {
	Type     string `yaml:"type" env-default:"postgres"`
	Name     string `yaml:"name"`
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
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

	jwtSecret := os.Getenv("JWT_SECRET")
	jwtTTLstr := os.Getenv("JWT_TTL")

	addr := os.Getenv("DATABASE_ADDRESS")
	pwd := os.Getenv("DATABASE_PASSWORD")
	dbUser := os.Getenv("DATABASE_USERNAME")
	dbName := os.Getenv("DATABASE_NAME")

	// check if file exists
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var config Config

	fmt.Println(configPath)

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	jwtTTL, err := time.ParseDuration(jwtTTLstr)
	if err != nil {
		log.Fatalf("error parse duration: %s", err)
	}

	config.JwtSecret = jwtSecret
	config.JwtTTL = jwtTTL

	config.DataBase.Address = addr
	config.DataBase.Username = dbUser
	config.DataBase.Password = pwd
	config.DataBase.Name = dbName

	return &config
}
