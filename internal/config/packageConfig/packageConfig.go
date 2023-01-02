package packageConfig

import (
	"encoding/json"
	"os"
)

type PackageConfig struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Architecture string `json:"architecture"`
	Files        []File `json:"files"`
}

type File struct {
	FilePath    string `json:"filePath"`
	InstallPath string `json:"installPath"`
	Static      bool   `json:"static"`
	RunScript   bool   `json:"runScript"`
	Script      string `json:"script"`
}

// Load : Loads a deb-studio application version configuration file
func Load(path string) (*PackageConfig, error) {
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
	av := &PackageConfig{}
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

// Save : Saves a deb-studio application version configuration file
func (av *PackageConfig) Save(path string) error {
	// Open config file
	configFile, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Create JSON from config object
	data, err := json.MarshalIndent(av, "", "\t")
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
