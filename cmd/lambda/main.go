package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"trackcoro/config"
	"trackcoro/controller"
	"trackcoro/database"
	"trackcoro/notify"
	"trackcoro/objectstorage"
	"trackcoro/router"
)

var initialized = false
var ginLambda *ginadapter.GinLambda

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !initialized {
		controller.InitializeControllers()
		r := router.InitializeRouter()
		ginLambda = ginadapter.New(r)
		initialized = true
	}
	return ginLambda.Proxy(req)
}

func main() {
	config.LoadConfig()
	database.ConnectToDB()
	defer database.DB.Close()
	database.MigrateSchema()
	notify.InitializeFirebase()
	objectstorage.InitializeS3Session()

	lambda.Start(Handler)
}
