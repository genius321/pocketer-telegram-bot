package config

import "github.com/spf13/viper"

type Config struct {
	Token       string
	ConsumerKey string
	RedirectURL string
	DBPath      string `mapstructure:"db_file"`

	Messages Messages
}

type Messages struct {
	Responses
	Errors
}

type Responses struct {
	StartMessage      string `mapstructure:"start_message"`
	SavedSuccessfully string `mapstructure:"saved_successfully"`
}

type Errors struct {
	Default           string `mapstructure:"default"`
	InvalidURL        string `mapstructure:"invalid_url"`
	UnableToSave      string `mapstructure:"unable_to_save"`
	AlreadyAuthorized string `mapstructure:"already_authorized"`
	UnknownCommand    string `mapstructure:"unknown_command"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config

	// for db_file
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.errors", &cfg.Messages.Errors); err != nil {
		return nil, err
	}

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
