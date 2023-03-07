package config

import (
	"flag"
	"os"
	"strings"
)

func (cfg *Config) envEnvType() {
	val, ok := os.LookupEnv("ENV")
	if ok {
		cfg.Env = AppMode(strings.ToLower(val))
	}
}

func (cfg *Config) flagEnvType() {
	val := flag.String("e", string(cfg.Env), "prod|debug")
	if val != nil {
		cfg.Env = AppMode(strings.ToLower(*val))
	}
}
