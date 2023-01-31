package config

import "errors"

/*
- адрес и порт запуска сервиса: переменная окружения ОС `RUN_ADDRESS` или флаг `-a`
- адрес подключения к базе данных: переменная окружения ОС `DATABASE_URI` или флаг `-d`
- адрес системы расчёта начислений: переменная окружения ОС `ACCRUAL_SYSTEM_ADDRESS` или флаг `-r`
*/

type Config struct {
	RunAddress           string // ОС `RUN_ADDRESS` или флаг `-a`
	DatabaseURI          string // ОС `DATABASE_URI` или флаг `-d`
	AccrualSystemAddress string // ОС `ACCRUAL_SYSTEM_ADDRESS` или флаг `-r`
}

func NewConfig() (Config, error) {
	cfg := Config{
		RunAddress:           "",
		DatabaseURI:          "",
		AccrualSystemAddress: "",
	}

	return cfg, errors.New("not implemented")
}
