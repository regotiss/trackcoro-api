package config

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"trackcoro/constants"
)

var Config *EnvConfig

type EnvConfig struct {
	FileServerURL      string
	DBConnectionString string
	FirebasePrivateKey string
}

func LoadConfig() {
	dbConnection, dbErr := getEnvVariable(constants.DBConnectionString)
	firebaseKey, fireBaseErr := getEnvVariable(constants.FirebaseKey)
	fileServerURL, fileServerErr := getEnvVariable(constants.FileServerURL)
	if dbErr != nil || fireBaseErr != nil || fileServerErr != nil {
		logrus.Panic("Could not find required environment variable(s)")
	}
	Config = &EnvConfig{
		FileServerURL:      fileServerURL,
		DBConnectionString: dbConnection,
		FirebasePrivateKey: firebaseKey,
	}
}

func getEnvVariable(key string) (string, error) {
	value := os.Getenv(key)
	if value == constants.Empty {
		return value, errors.New(constants.EnvVariableNotFoundError)
	}
	return value, nil
}
