package helpers

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var loadPath = "/app/.env"

func ErrorHelper(err error, msg string) error {
	return fmt.Errorf("%s:%v", msg, err)
}

func Token() (string, error) {
	err := godotenv.Load(loadPath)
	if err != nil {
		return "", err
	}
	token := os.Getenv("SECRET")

	if token == "" {
		return "", ErrorHelper(fmt.Errorf("token is empty"), "error")
	}

	return token, nil
}
