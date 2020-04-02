package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	r := InitializeRouter()
	err := r.Run()
	if err != nil {
		logrus.Error("Could not start server", err)
	}
}
