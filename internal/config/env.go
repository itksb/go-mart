package config

// UseOsEnv - apply environment variables
func (cfg *Config) UseOsEnv() {
	cfg.envRunAddress()
	cfg.envDatabaseURI()
	cfg.envAccrualSystemAddress()
}
