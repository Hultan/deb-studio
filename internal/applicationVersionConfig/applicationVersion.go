package applicationVersionConfig

import (
	"encoding/json"
	"os"
)

type ApplicationVersion struct {
	Control     ApplicationVersionControl `json:"control"`
	PreInstall  string                    `json:"preInstall"`
	Files       []File                    `json:"files"`
	PostInstall string                    `json:"postInstall"`
}

type ApplicationVersionControl struct {
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
}

type File struct {
	Path         string `json:"path"`
	Copy         bool   `json:"copy"`
	CopyFromPath string `json:"copyFromPath"`
}

// Load : Loads a deb-studio application version configuration file
func Load(path string) (*ApplicationVersion, error) {
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
	av := &ApplicationVersion{}
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
func (av *ApplicationVersion) Save(path string) error {
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
