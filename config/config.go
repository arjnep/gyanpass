package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ServerConfiguration struct {
	Port      string
	JWTSecret string
	JWTExpiry int
	Timeout   int
	Mode      string
	Version   string
}

type DatabaseConfiguration struct {
	DBName       string
	Username     string
	Password     string
	Host         string
	Port         string
	SSLMode      string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
}

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
}

var config *Configuration

func LoadConfig() {
	if os.Getenv("SERVER_MODE") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error Loading Environment Variables: %v", err)
			return
		}
	}

	ctxTimeout, _ := strconv.Atoi(os.Getenv("SERVER_TIMEOUT"))
	dbMaxLifetime, _ := strconv.Atoi(os.Getenv("DATABASE_MAX_LIFETIME"))
	dbMaxOpenConns, _ := strconv.Atoi(os.Getenv("DATABASE_MAX_OPEN_CONNS"))
	dbMaxIdleConns, _ := strconv.Atoi(os.Getenv("DATABASE_MAX_IDLE_CONNS"))
	jWTExpiry, _ := strconv.Atoi(os.Getenv("SERVER_JWTExpiry"))

	cfg := &Configuration{
		Server: ServerConfiguration{
			Port:      os.Getenv("SERVER_PORT"),
			JWTSecret: os.Getenv("SERVER_JWTSECRET"),
			JWTExpiry: jWTExpiry,
			Timeout:   ctxTimeout,
			Mode:      os.Getenv("SERVER_MODE"),
			Version:   os.Getenv("SERVER_VERSION"),
		},
		Database: DatabaseConfiguration{
			DBName:       os.Getenv("DATABASE_DBNAME"),
			Username:     os.Getenv("DATABASE_USERNAME"),
			Password:     os.Getenv("DATABASE_PASSWORD"),
			Host:         os.Getenv("DATABASE_HOST"),
			Port:         os.Getenv("DATABASE_PORT"),
			SSLMode:      os.Getenv("DATABASE_SSLMODE"),
			MaxLifetime:  dbMaxLifetime,
			MaxOpenConns: dbMaxOpenConns,
			MaxIdleConns: dbMaxIdleConns,
		},
	}

	config = cfg
	log.Println("Environment Variables Loaded!")
}

func GetConfig() *Configuration {
	return config
}
