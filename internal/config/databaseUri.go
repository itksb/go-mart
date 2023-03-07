package config

import (
	"flag"
	"os"
)

func (cfg *Config) envDatabaseURI() {
	dsn, ok := os.LookupEnv("DATABASE_URI")
	if ok {
		cfg.DatabaseURI = dsn
	}
}

func (cfg *Config) flagDatabaseURI() {
	dsn := flag.String("d", cfg.DatabaseURI, "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable")
	if dsn != nil {
		cfg.DatabaseURI = *dsn
	}
}
