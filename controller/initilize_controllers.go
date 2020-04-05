package controller

import (
	"trackcoro/admin"
	"trackcoro/database"
	"trackcoro/quarantine"
	"trackcoro/so"
)

var (
	QuarantineController quarantine.Controller
	AdminController      admin.Controller
	SOController         so.Controller
)

func InitializeControllers() {
	quarantineRepo := quarantine.NewRepository(database.DB)
	quarantineService := quarantine.NewService(quarantineRepo)
	QuarantineController = quarantine.NewController(quarantineService)

	adminRepo := admin.NewRepository(database.DB)
	adminService := admin.NewService(adminRepo)
	AdminController = admin.NewController(adminService)

	soRepo := so.NewRepository(database.DB)
	soService := so.NewService(soRepo)
	SOController = so.NewController(soService)
}
