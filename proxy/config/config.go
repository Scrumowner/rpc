package config

type Config struct {
}

type Server struct {
}

type Database struct {
	Port int
}

type Cache struct {
	Address  string `yaml:"address"`
	Password string `json:"-" yaml:"password"`
	Port     string `yaml:"port"`
}
