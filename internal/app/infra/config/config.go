package config

import "github.com/ryrden/rinha-de-backend-go/pkg/env"

type Config struct {
	Database Database

	// TODO: Config this later
	Cache     Cache
	Server    Server
	Profiling Profiling
}

type Database struct {
	Host               string
	Port               string
	User               string
	Password           string
	Name               string
	Max_db_connections string
	Min_db_connections string
}

type Cache struct {
	Host string
	Port string
}

type Server struct {
	Port     string
	UseSonic bool
	Prefork  bool
}

type Profiling struct {
	Enabled bool
	CPU     string
	Mem     string
}

func NewConfig() *Config {
	return &Config{
		Database: Database{
			Host:               env.GetEnvOrDie("DB_HOST"),
			Port:               env.GetEnvOrDie("DB_PORT"),
			User:               env.GetEnvOrDie("DB_USER"),
			Password:           env.GetEnvOrDie("DB_PASSWORD"),
			Name:               env.GetEnvOrDie("DB_NAME"),
			Max_db_connections: env.GetEnvOrDie("MAX_DB_CONNECTIONS"),
			Min_db_connections: env.GetEnvOrDie("MIN_DB_CONNECTIONS"),
		},
		Cache: Cache{
			Host: env.GetEnvOrDie("CACHE_HOST"),
			Port: env.GetEnvOrDie("CACHE_PORT"),
		},
		Server: Server{
			Port:     env.GetEnvOrDie("SERVER_PORT"),
			UseSonic: env.GetEnvOrDie("ENABLE_SONIC_JSON") == "1",
			Prefork:  env.GetEnvOrDie("ENABLE_PREFORK") == "1",
		},
		Profiling: Profiling{
			Enabled: env.GetEnvOrDie("ENABLE_PROFILING") == "1",
			CPU:     env.GetEnvOrDie("CPU_PROFILE"),
			Mem:     env.GetEnvOrDie("MEM_PROFILE"),
		},
	}
}
