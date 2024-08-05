package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		App      `yaml:"app"`
		Server   `yaml:"http"`
		Oauth    `yaml:"oauth"`
		JwtAuth  `yaml:"jwtAuth"`
		Log      `yaml:"logger"`
		Database `yaml:"postgres"`
	}

	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}
	Server struct {
		Port string `yaml:"port"`
	}

	Oauth struct {
		ClientID         string `yaml:"clientId"`
		ClientSecret     string `yaml:"clientSecret"`
		RedirectURL      string `yaml:"redirectUrl"`
		OauthStateString string `yaml:"stateString"`
	}
	JwtAuth struct {
		Key string `yaml:"key"`
	}

	Log struct {
		Level string `yaml:"log_level"`
	}
	Database struct {
		// Host         string `yaml:"host"`
		// Port         int    `yaml:"port"`
		// User         string `yaml:"user"`
		// Password     string `yaml:"password"`
		// DatabaseName string `yaml:"db_name"`
		ConnString string `yaml:"conn_string"`
	}
)

func New(cfgPath string) (*Config, error) {
	cfg := &Config{}

	file, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
