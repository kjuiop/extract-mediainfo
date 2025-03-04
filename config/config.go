package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Logger  Logger  `json:"log"`
	Watcher Watcher `json:"watcher"`
}

type Watcher struct {
	CheckIntervalSec int      `json:"check_interval_sec"`
	Targets          []Target `json:"targets"`
}

type Logger struct {
	Path        string `json:"path"`
	Level       string `json:"level"`
	PrintStdOut bool   `json:"print_std_out"`
}

type Target struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

func ReadConfig() (*Config, error) {
	c := new(Config)
	configFile, err := os.Open("conf/conf.json")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := configFile.Close(); err != nil {
			log.Fatalf("Read Config file Fail, err : %s", err.Error())
		}
	}()

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&c); err != nil {
		return nil, err
	}

	if err := c.CheckValid(); err != nil {
		return nil, fmt.Errorf("check config, err: %w", err)
	}

	return c, nil
}

func (c *Config) CheckValid() error {
	return nil
}
