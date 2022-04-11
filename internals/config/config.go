package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DBHost   string `json:"host"`
	DBName   string `json:"db"`
}

func GetConfig(filepath string) *Config {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &Config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
