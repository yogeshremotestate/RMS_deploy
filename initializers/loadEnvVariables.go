package initializers

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var ENV struct {
	PORT           string
	DB_URL         string
	SECRET         string
	MIGRATIONS_URL string
}

//func LoadEnvVariables() {
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}
//
//	ENV.PORT = os.Getenv("PORT")
//	ENV.DB_URL = os.Getenv("DB_URL")
//	ENV.SECRET = os.Getenv("SECRET")
//
//}

var UserString = "user"

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ENV.PORT = os.Getenv("PORT")
	ENV.DB_URL = os.Getenv("DB_URL")
	ENV.SECRET = os.Getenv("SECRET")
	ENV.MIGRATIONS_URL = os.Getenv("MIGRATIONS_URL")
}

//func LoadEnvVariables() {
//	// Directly fetch environment variables from ECS
//	ENV.PORT = os.Getenv("PORT")
//	if ENV.PORT == "" {
//		log.Fatal("PORT environment variable not set")
//	}
//
//	ENV.DB_URL = os.Getenv("DB_URL")
//	if ENV.DB_URL == "" {
//		log.Fatal("DB_URL environment variable not set")
//	}
//
//	ENV.SECRET = os.Getenv("SECRET")
//	if ENV.SECRET == "" {
//		log.Fatal("SECRET environment variable not set")
//	}
//
//	ENV.MIGRATIONS_URL = os.Getenv("MIGRATIONS_URL")
//	if ENV.MIGRATIONS_URL == "" {
//		log.Fatal("MIGRATIONS_URL environment variable not set")
//	}
//}
