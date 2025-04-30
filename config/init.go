package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	DefaultPhoto string
)

type Config struct {
	Database     Database
	GeneralPhoto General
}

type Database struct {
	dbHost string
	dbPort int
	dbName string
	dbPass string
	dbUser string
	SSL    string
}

type General struct {
	DefaultPhoto string
}

func EnvFile() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error load .env data")
	}
	port, _ := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	return &Config{
		Database: Database{
			dbHost: os.Getenv("DATABASE_HOST"),
			dbPort: port,
			dbUser: os.Getenv("DATABASE_USER"),
			dbPass: os.Getenv("DATABASE_PASS"),
			dbName: os.Getenv("DATABASE_NAME"),
			SSL:    os.Getenv("DATABASE_SSL"),
		},
		GeneralPhoto: General{
			DefaultPhoto: os.Getenv("DEFAULT_PP"),
		},
	}
}
