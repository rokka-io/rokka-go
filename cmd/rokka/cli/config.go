package cli

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// Config stores the APIKey and ImageHost in a file to be used later for authenticating against rokka without having to pass an APIKey/ImageHost as a flag.
type Config struct {
	APIKey    string `json:"apiKey"`
	ImageHost string `json:"imageHost"`
}

var configPath string

// SetConfigPath adjusts the global variable configPath which tells where to look for the configuration file.
func SetConfigPath(p string) {
	configPath = p
}

// GetConfig returns the stored config if set. If not it will return an empty configuration object.
func GetConfig() (Config, error) {
	config := Config{}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil
	}

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	if err = json.Unmarshal(content, &config); err != nil {
		return config, err
	}

	return config, nil
}

// SaveConfig stores the current configuration in configPath.
func SaveConfig(c Config) error {
	dir := path.Dir(configPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	d, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f, err := os.Create(configPath)
	if err != nil {
		return err
	}

	_, err = f.Write(d)
	if err != nil {
		return err
	}

	return nil
}
