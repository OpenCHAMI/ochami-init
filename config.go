package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Networks []struct {
		Name   string `yaml:"name"`
		Subnet string `yaml:"subnet"`
	} `yaml:"networks"`
	Databases []struct {
		Name  string `yaml:"name"`
		Users []struct {
			Name     string `yaml:"name"`
			Password string `yaml:"password"`
		} `yaml:"users"`
	} `yaml:"databases"`
}

func readConfig(filename string) (*Config, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func writeConfig(filename string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
