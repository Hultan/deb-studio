package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Control struct {
		Package       string `json:"package"`
		Source        string `json:"source"`
		Version       string `json:"version"`
		Section       string `json:"section"`
		Priority      string `json:"priority"`
		Architecture  string `json:"architecture"`
		Essential     bool   `json:"essential"`
		Depends       string `json:"depends"`
		InstalledSize string `json:"installed-size"`
		Maintainer    string `json:"maintainer"`
		Description   string `json:"description"`
		Homepage      string `json:"homepage"`
		BuiltUsing    string `json:"built-using"`
	} `json:"control"`
}

func Load(path string) (*Config, error) {
	// Make sure the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("settings file is missing (%s)", path))
	}

	// Open config file
	configFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to open config")
	}

	// Parse the JSON document
	var config Config
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return nil, err
	}

	err = configFile.Close()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Save : Saves a SoftTube configuration file
func (c *Config) Save(path string) error {
	// Open config file
	configFile, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Create JSON from config object
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	// Write the data
	_, err = configFile.Write(data)
	if err != nil {
		return err
	}

	err = configFile.Close()
	if err != nil {
		return err
	}

	return nil
}
