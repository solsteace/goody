package internal

import (
	"log"
	"os"
	"strconv"
)

var (
	// The port the app be run on
	EnvPort uint

	// The url of the main database
	EnvDbUrl string

	// The Url of message queue
	EnvMqUrl string

	// The base path used for saving uploaded files
	EnvUploadSaveBasePath string
)

func loadEnv() {
	port, err := strconv.ParseUint(os.Getenv("PORT"), 10, 64)
	if err != nil {
		log.Fatalf("`PORT`: %v", err)
	} else if port := uint(port); port < 0 || port > 65535 {
		log.Fatal("`PORT` should be between 0 - 65535")
	}
	EnvPort = uint(port)

	EnvMqUrl = os.Getenv("MQ_URL")
	EnvDbUrl = os.Getenv("DB_URL")
	EnvUploadSaveBasePath = os.Getenv("UPLOAD_SAVE_BASE_PATH")
}
