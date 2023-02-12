package config

import "os"

type Config struct {
	DB_URL string
	PORT   string
}

func NewConfig() *Config {
	env := os.Getenv("ENV")
	if env == "PROD" {
		return &Config{
			DB_URL: os.Getenv("DB_URL"),
			PORT:   os.Getenv("PORT"),
		}
	} else {
		return &Config{
			DB_URL: "mongodb://root:password@localhost:27017",
			PORT:   "8080",
		}
	}
}
