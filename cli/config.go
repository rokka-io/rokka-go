package cli

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

type Config struct {
	APIKey       string `json:"apiKey"`
	APIAddress   string `json:"apiAddress"`
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

	p := getPath(usr)

	if _, err := os.Stat(p); os.IsNotExist(err) {
		return config, nil
	}

	content, err := ioutil.ReadFile(p)
	if err != nil {
		return config, err
	}

	if err = json.Unmarshal(content, &config); err != nil {
		return config, err
	}

	return config, nil
}
