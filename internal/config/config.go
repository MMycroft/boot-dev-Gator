// Package config holds config functions
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const _FILE_NAME = ".gatorconfig.json"

func Read() (Config, error) {
	var cfg Config
	directory, err := os.UserHomeDir()
	if err != nil {
		return cfg, err
	}

	file, err := os.Open(filepath.Join(directory, _FILE_NAME))
	if err != nil {
		return cfg, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (cfg *Config) SetUser(userName string) error {
	directory, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	readFile, err := os.Open(filepath.Join(directory, _FILE_NAME))
	if err != nil {
		return err
	}

	err = json.NewDecoder(readFile).Decode(&cfg)
	if err != nil {
		return err
	}
	readFile.Close()

	cfg.CurrentUserName = userName

	writeFile, err := os.Create(filepath.Join(directory, _FILE_NAME))
	if err != nil {
		return err
	}

	err = json.NewEncoder(writeFile).Encode(cfg)
	if err != nil {
		return err
	}
	writeFile.Close()

	return nil
}
