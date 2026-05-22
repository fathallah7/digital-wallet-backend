package config

import "os"

type Config struct {
	Port  string
	DBUrl string
}

func Load() Config {

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	return Config{
		Port: port,
		DBUrl: os.Getenv("DB_DSN"),
	}

}
