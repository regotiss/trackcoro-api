package admin

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"trackcoro/constants"
	"trackcoro/database/models"
)

type Service interface {
	Verify(mobileNumber string) bool
	Add() error
}

type service struct {
	repository Repository
}

func (s service) Verify(mobileNumber string) bool {
	return s.repository.IsExists(mobileNumber)
}

func (s service) Add() error {
	mobileNumber := os.Getenv(constants.AdminMobileNumber)
	if mobileNumber == constants.Empty {
		logrus.Error(constants.EnvVariableNotFound)
		return errors.New(constants.EnvVariableNotFound)
	}
	return s.repository.Add(models.Admin{MobileNumber:mobileNumber})
}

func NewService(repository Repository) Service {
	return service{repository}
}
