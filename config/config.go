package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var GEMINI_API_KEY string

func Init() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	GEMINI_API_KEY = os.Getenv("GEMINI_API_KEY")

	if GEMINI_API_KEY == "" {
		return errors.New("GEMINI_API_KEY environment variable is not set")
	}

	return nil
}
