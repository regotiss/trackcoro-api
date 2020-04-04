package controller

import (
	"trackcoro/database"
	"trackcoro/quarantine"
)

var (
	QuarantineController quarantine.Controller
)

func InitializeControllers() {
	repo := quarantine.NewRepository(database.DB)
	service := quarantine.NewService(repo)
	QuarantineController = quarantine.NewController(service)
}