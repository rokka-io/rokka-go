package cli

import (
	"encoding/json"
	"os"
	"os/user"
	"path"	
	"io/ioutil"
)

type Config struct {
	APIKey string `json:"apiKey"`
	Organization string `json:"organization"`
}

func getPath(u *user.User) string {
	return path.Join(u.HomeDir, ".rokka", "config")
}

func GetConfig() (Config, error) {
	config := Config{}

	usr, err := user.Current()
	if err != nil {
		return config, err
	}

	path := getPath(usr)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return config, nil
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	if err = json.Unmarshal(content, &config); err != nil {
		return config, err
	}

	return config, nil
}