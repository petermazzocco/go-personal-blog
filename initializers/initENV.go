package initializers

import (
	"fmt"

	"github.com/joho/godotenv"
)

func InitENV() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file, skipping..")
	}
}
