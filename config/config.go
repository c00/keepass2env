package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Entries      []Entry `yaml:"entries"`
	DatabasePath string  `yaml:"databasePath,omitempty"`
	KeyfilePath  string  `yaml:"keyfilePath,omitempty"`
	// The environment variable to check for the password
	PasswordEnv string `yaml:"passwordEnv,omitempty"`
	// The location where the .env file will be created / updated
	OutputPath string `yaml:"outputPath,omitempty"`
}

type Entry struct {
	// The Environment Variable name to set
	EnvName string `yaml:"envName"`
	// The path inside the Keepass Database
	KeepassPath string `yaml:"keepassPath"`
}

func FromFile(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("cannot open config: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("cannot decode yaml: %w", err)
	}

	return cfg, nil
}
