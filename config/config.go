package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	HttpPort string   `json:"httpport"`
	Database Database `json:"db"`
}

type Database struct {
	DbHost string `json:"host"`
	DbPort string `json:"port"`
	DbUser string `json:"user"`
	DbPass string `json:"password"`
	DbName string `json:"database"`
}

func NewConfiguration(file string) Config {
	config := Config{}

	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return config
}

func (c *Config) OverrideFromEnvironment() {
	if v := os.Getenv("APP_PORT"); v != "" {
		c.HttpPort = v
	}
	if v := os.Getenv("DB_HOST"); v != "" {
		c.Database.DbHost = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		c.Database.DbPort = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		c.Database.DbUser = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		c.Database.DbPass = v
	}
	if v := os.Getenv("DB_DBNAME"); v != "" {
		c.Database.DbName = v
	}
}
