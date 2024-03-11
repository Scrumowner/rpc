package config

import "os"

type Config struct {
	*ServConfig
	*HugoAddr
	*GeoAddr
	*AuthAddr
	*UserAddr
}
type ServConfig struct {
	Listen string
	Port   string
}

type HugoAddr struct {
	Port string
}
type GeoAddr struct {
	Host string
	Port string
}
type AuthAddr struct {
	Host string
	Port string
}
type UserAddr struct {
	Host string
	Port string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load() {
	c.ServConfig = &ServConfig{
		Listen: os.Getenv("USER_LISTEN"),
		Port:   os.Getenv("PORXY_PORT"),
	}
	c.HugoAddr = &HugoAddr{
		Port: os.Getenv("HUGO_PORT"),
	}
	c.GeoAddr = &GeoAddr{
		Host: os.Getenv("GEO_HOST"),
		Port: os.Getenv("GEO_PORT"),
	}
	c.AuthAddr = &AuthAddr{
		Host: os.Getenv("AUTH_HOST"),
		Port: os.Getenv("AUTH_PORT"),
	}
	c.UserAddr = &UserAddr{
		Host: os.Getenv("USER_HOST"),
		Port: os.Getenv("USER_PORT"),
	}
}
