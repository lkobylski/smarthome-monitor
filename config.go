package main

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Broker            string   `yaml:"broker"`
	ClientID          string   `yaml:"client_id"`
	Topics            []string `yaml:"topics"`
	DisconnectTimeout int      `yaml:"disconnect_timeout"`
	Devices           []Device `yaml:"devices"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
