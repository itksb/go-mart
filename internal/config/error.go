package config

import "errors"

var ErrConfigWrongEnvValue = errors.New("config ENV wrong value")
var ErrConfigAppHostIsEmpty = errors.New("wrong configuration. AppHost is empty.")
var ErrConfigDatabaseURIEmpty = errors.New("wrong configuration. DatabaseURI is empty.")
