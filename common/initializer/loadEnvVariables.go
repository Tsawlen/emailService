package initializer

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables(doneChannel chan bool) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	doneChannel <- true
}
