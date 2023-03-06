package config

import (
	"fmt"
	"time"
)

/*
- адрес и порт запуска сервиса: переменная окружения ОС `RUN_ADDRESS` или флаг `-a`
- адрес подключения к базе данных: переменная окружения ОС `DATABASE_URI` или флаг `-d`
- адрес системы расчёта начислений: переменная окружения ОС `ACCRUAL_SYSTEM_ADDRESS` или флаг `-r`
*/

type Config struct {
	AppHost                  string // ОС `RUN_ADDRESS` или флаг `-a`, e.g.: localhost:8080
	AppPort                  int    // ОС `RUN_ADDRESS` или флаг `-a`
	DatabaseURI              string // ОС `DATABASE_URI` или флаг `-d`
	AccrualSystemAddress     string // ОС `ACCRUAL_SYSTEM_ADDRESS` или флаг `-r`
	AppSecret                string
	Env                      AppMode
	GracefulShutdownInterval time.Duration
}

type AppMode string

const (
	Debug AppMode = "debug"
	Prod  AppMode = "prod"
)

func NewConfig() (Config, error) {
	cfg := Config{
		AppHost:                  "",
		AppPort:                  8080,
		DatabaseURI:              "",
		AccrualSystemAddress:     "",
		AppSecret:                "demo.demo.demo",
		Env:                      Prod,
		GracefulShutdownInterval: 5, //seconds
	}

	return cfg, nil
}

func (cfg *Config) GetFullAddr() string {
	return fmt.Sprintf("%s:%d", cfg.AppHost, cfg.AppPort)
}

// UseOsEnv - apply environment variables
func (cfg *Config) UseOsEnv() {
	cfg.envRunAddress()
	cfg.envDatabaseURI()
	cfg.envAccrualSystemAddress()
	cfg.envEnvType()
}

// UseFlags - scan flags
func (cfg *Config) UseFlags() {
	cfg.flagRunAddress()
	cfg.flagDatabaseURI()
	cfg.flagAccrualSystemAddress()
	cfg.flagEnvType()
}
