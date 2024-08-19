package database

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load("C:/Users/hp/goLang/goLang-Project-2024/goProject-2024/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
