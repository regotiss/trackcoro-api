package main

import (
	"github.com/sirupsen/logrus"
	"trackcoro/controller"
	"trackcoro/database"
	"trackcoro/notify"
	"trackcoro/objectstorage"
	"trackcoro/router"
)

func main() {
	database.ConnectToDB()
	defer database.DB.Close()
	database.MigrateSchema()
	objectstorage.InitializeS3Session()
	controller.InitializeControllers()
	notify.InitializeFirebase()
	r := router.InitializeRouter()
	err := r.Run()
	if err != nil {
		logrus.Error("Could not start server", err)
	}

}
