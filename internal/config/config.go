package config

import (
	"encoding/json"
	"flag"
	"io"
	"os"
)

type Config struct {
	Listen      string `json:"listen"`
	Production  bool   `json:"production"`
	Backend     string `json:"backend"`
	S3Endpoint  string `json:"s3_endpoint"`
	S3AccessKey string `json:"s3_access_key"`
	S3SecretKey string `json:"s3_secret_key"`
	PostgresURL string `json:"postgres_url"`
}

func New() (Config, error) {
	config := flag.String("config", "configs/jobdoe.json", "Path to configuration file")
	flag.Parse()

	result := Config{}

	configFile, err := os.Open(*config)
	if err != nil {
		return result, err
	}
	defer configFile.Close()

	configRaw, err := io.ReadAll(configFile)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(configRaw, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}
