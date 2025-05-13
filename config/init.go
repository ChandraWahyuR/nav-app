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
	SMTP         SMTP
	Gmaps        GMAPS
	URL_Server   string
}

type Database struct {
	dbHost string
	dbPort int
	dbName string
	dbPass string
	dbUser string
	SSL    string
	cert   string
}

type General struct {
	DefaultPhoto string
}

type SMTP struct {
	SMTP_HOST          string
	SMTP_PORT          int
	SMTP_EMAIL_ADDRESS string
	SMTP_TOKEN_EMAIL   string
}

type GMAPS struct {
	GMAPS_API_KEY string
}

func EnvFile() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error load .env data")
	}
	port, _ := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	portSMTP, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	return &Config{
		Database: Database{
			dbHost: os.Getenv("DATABASE_HOST"),
			dbPort: port,
			dbUser: os.Getenv("DATABASE_USER"),
			dbPass: os.Getenv("DATABASE_PASS"),
			dbName: os.Getenv("DATABASE_NAME"),
			SSL:    os.Getenv("DATABASE_SSL"),
			cert:   os.Getenv("DATABASE_CERT"),
		},
		GeneralPhoto: General{
			DefaultPhoto: os.Getenv("DEFAULT_PP"),
		},
		SMTP: SMTP{
			SMTP_HOST:          os.Getenv("SMTP_HOST"),
			SMTP_PORT:          portSMTP,
			SMTP_EMAIL_ADDRESS: os.Getenv("SMTP_EMAIL_ADDRESS"),
			SMTP_TOKEN_EMAIL:   os.Getenv("SMTP_TOKEN_EMAIL"),
		},
		Gmaps: GMAPS{
			GMAPS_API_KEY: os.Getenv("GMAPS_API_KEY"),
		},
		URL_Server: os.Getenv("ENDPOINT_SERVER"),
	}
}
