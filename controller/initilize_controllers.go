package controller

import (
	"trackcoro/admin"
	"trackcoro/database"
	"trackcoro/quarantine"
)

var (
	QuarantineController quarantine.Controller
	AdminController      admin.Controller
)

func InitializeControllers() {
	quarantineRepo := quarantine.NewRepository(database.DB)
	quarantineService := quarantine.NewService(quarantineRepo)
	QuarantineController = quarantine.NewController(quarantineService)

	adminRepo := admin.NewRepository(database.DB)
	adminService := admin.NewService(adminRepo)
	AdminController = admin.NewController(adminService)
}
