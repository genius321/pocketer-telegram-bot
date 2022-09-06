package config

import "github.com/spf13/viper"

type Config struct {
	Token       string
	ConsumerKey string
	RedirectURL string
}

func Init() (*Config, error) {
	var cfg Config

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseEnv(cfg *Config) error {
	if err := viper.BindEnv("token"); err != nil {
		return err
	}

	if err := viper.BindEnv("consumer_key"); err != nil {
		return err
	}

	if err := viper.BindEnv("redirect_url"); err != nil {
		return err
	}

	cfg.Token = viper.GetString("token")
	cfg.ConsumerKey = viper.GetString("consumer_key")
	cfg.RedirectURL = viper.GetString("redirect_url")

	return nil
}
