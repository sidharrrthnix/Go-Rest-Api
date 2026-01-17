package utils

import (
	"fmt"
	"log"
	"os"
)

var errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

func ErrorHandler(err error, message string) error {
	errorLogger.Println(message, err)
	return fmt.Errorf("%s: %w", message, err)
}
