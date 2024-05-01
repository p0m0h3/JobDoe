package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Listen       string   `json:"listen"`
	Mode         string   `json:"mode"`
	Podman       string   `json:"podman"`
	Key          string   `json:"key"`
	RegistryAuth string   `json:"regauth"`
	Proxies      []string `json:"proxies"`
}

var config = &Config{}

func GetConfig() (*Config, error) {
	if config.Listen != "" {
		return config, nil
	}

	file, err := os.ReadFile("config.json")
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(file, &config)
	return config, err

}
