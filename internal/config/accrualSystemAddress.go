package config

import (
	"flag"
	"os"
)

func (cfg *Config) envAccrualSystemAddress() {
	addr, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if ok {
		cfg.AccrualSystemAddress = addr
	}
}

func (cfg *Config) flagAccrualSystemAddress() {
	addr := flag.String("r", cfg.AppHost, "ACCRUAL_SYSTEM_ADDRESS")
	if addr != nil {
		cfg.AccrualSystemAddress = *addr
	}
}
