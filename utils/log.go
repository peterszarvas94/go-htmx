package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Loglevel int

const (
	INFO Loglevel = iota
	WARNING
	ERROR
	FATAL
)

const maxLoglevel = int(FATAL)

// Returns the string representation of the loglevel
func llToString(ll Loglevel) string {
	switch ll {
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "INFO"
	}
}

func llToEnum(ll string) Loglevel {
	switch ll {
	case "INFO":
		return INFO
	case "WARNING":
		return WARNING
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		return INFO
	}
}

// Logs into log.txt in the following format:
// LEVEL location YYYY-MM-DD HH:MM:SS - message
func Log(messageLl Loglevel, location string, message string) {
	fileName := "log.txt"
	now := time.Now().Format("2006-01-02 15:04:05")

	file, fileErr := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fileErr != nil {
		logger := log.New(os.Stdout, "", 0)
		logger.Printf("ERROR log/file %s - %s", now, fileErr.Error())
	}
	defer file.Close()

	messageLlStr := llToString(messageLl)

	formattedMessage := fmt.Sprintf("%s %s %s - %s", messageLlStr, location, now, message)
	logger := log.New(file, "", 0)

	ll := GetLogLevel()

	if messageLl < ll {
		return
	}

	if messageLl == FATAL {
		logger.Fatal(formattedMessage)
	} else {
		logger.Println(formattedMessage)
	}
}

// Returns the loggin level from the LOG_LEVEL environment variable, defaulting to INFO
func GetLogLevel() Loglevel {
	llStr, ok := os.LookupEnv("LOG_LEVEL")
	if !ok || llStr == "" {
		return INFO
	}

	ll := llToEnum(llStr)
	llInt := int(ll)

	if llInt < 0 || llInt > maxLoglevel {
		return INFO
	}

	return ll
}
