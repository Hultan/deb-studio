package projectConfig

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

type ProjectConfig struct {
	Name                  string    `json:"name"`
	LatestVersion         string    `json:"latestVersion"`
	CurrentPackageId      uuid.UUID `json:"currentPackage"`
	ShowOnlyLatestVersion bool      `json:"showOnlyLatestVersion"`
}

// Load : Loads a deb-studio config file
func Load(path string) (*ProjectConfig, error) {
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
	av := &ProjectConfig{}
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
func (c *ProjectConfig) Save(path string) error {
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
