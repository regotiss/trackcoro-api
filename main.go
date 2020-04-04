package main

import (
	"github.com/sirupsen/logrus"
	"trackcoro/database"
)

func main() {
	database.ConnectToDB()
	defer database.DB.Close()
	database.MigrateSchema()

	InitializeControllers()
	r := InitializeRouter()
	err := r.Run()
	if err != nil {
		logrus.Error("Could not start server", err)
	}

}
