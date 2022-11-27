package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	WorkingFolders []WorkingFolder `json:"workingFolders"`
}

type WorkingFolder struct {
	Path string `json:"path"`
}

// Load : Loads a deb-studio config file
func Load(path string) (*Config, error) {
	// Make sure the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	// Open config file
	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Parse the JSON document
	av := &Config{}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(av)
	if err != nil {
		return nil, err
	}

	err = configFile.Close()
	if err != nil {
		return nil, err
	}

	return av, nil
}

// Save : Saves a deb-studio config file
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
