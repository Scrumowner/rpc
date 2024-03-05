package config

import "os"

type Config struct {
	Listen string
	Port   string
	Secret string
	*UserService
}

type UserService struct {
	Host string
	Port string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load() {
	c.Listen = os.Getenv("AUTH_LISTEN")
	c.Port = os.Getenv("AUTH_PORT")
	c.Secret = os.Getenv("AUTH_SECRET")
	c.UserService = &UserService{
		Host: os.Getenv("USER_HOST"),
		Port: os.Getenv("USER_PORT"),
	}

}
