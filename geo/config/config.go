package config

import "os"

type Config struct {
	Listen string
	Port   string
	TokenA string
	TokenX string
	*DbConfig
	*CacheConfig
}
type DbConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Dbname   string
}
type CacheConfig struct {
	Host string
	Port string
}

func NewConfig() *Config {
	cfg := Config{}
	return &cfg
}

func (c *Config) Load() {
	c.Listen = os.Getenv("GEO_LISTEN")
	c.Port = os.Getenv("GEO_PORT")
	c.TokenA = os.Getenv("USER_TOKEN_A")
	c.TokenX = os.Getenv("USER_TOKEN_X")
	c.DbConfig = &DbConfig{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Dbname:   os.Getenv("POSTGRES_DB"),
	}
	c.CacheConfig = &CacheConfig{
		Host: os.Getenv("REDIS_HOST"),
		Port: os.Getenv("REDIS_PORT"),
	}

}
