package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Token           string
	DBConnURL       string
	RedisAddr       string
	DBAttendanceURL string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Token:           os.Getenv("TOKEN"),
		DBConnURL:       os.Getenv("DATABASE_URL"),
		RedisAddr:       os.Getenv("REDIS_ADDR"),
		DBAttendanceURL: os.Getenv("DATABASE_ATTENDANCE_URL"),
	}
}
