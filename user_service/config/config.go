package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	Dbname   string
	*ServConfig
}

type ServConfig struct {
	Listen string
	Port   string
}

func NewcConfig() *Config {
	return &Config{}
}

func (c *Config) Load() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("User_service can't read .env ")
	}
	c.User = os.Getenv("POSTGRES_USER")
	c.Password = os.Getenv("POSTGRES_PASSWORD")
	c.Host = os.Getenv("POSTGRES_HOST")
	c.Port = os.Getenv("POSTGRES_PORT")
	c.Dbname = os.Getenv("POSTGRES_DB")
	c.ServConfig = &ServConfig{
		Listen: os.Getenv("USER_LISTEN"),
		Port:   os.Getenv("USER_PORT"),
	}

}
