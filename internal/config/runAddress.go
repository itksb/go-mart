package config

import (
	"flag"
	"log"
	"os"
)

func verifyAndSetup(runAddress string, cfg *Config) {
	host, port, err := extractRunAddress(runAddress)
	if err != nil {
		log.Panic(err)
	}
	if host != "" {
		cfg.AppHost = host
	}
	if port != 0 {
		cfg.AppPort = port
	}

}

func (cfg *Config) envRunAddress() {
	runAddress, ok := os.LookupEnv("RUN_ADDRESS")
	if ok {
		verifyAndSetup(runAddress, cfg)
	}

}

func (cfg *Config) flagRunAddress() {
	runAddress := flag.String("a", cfg.AppHost, "RUN_ADDRESS")
	if runAddress != nil {
		verifyAndSetup(*runAddress, cfg)
	}
}
