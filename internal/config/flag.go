package config

// UseFlags - scan flags
func (cfg *Config) UseFlags() {
	cfg.flagRunAddress()
	cfg.flagDatabaseURI()
	cfg.flagAccrualSystemAddress()
}
