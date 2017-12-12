package cli

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	APIKey string `json:"apiKey"`
}

var configPath string

func SetConfigPath(p string) {
	configPath = p
}

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
