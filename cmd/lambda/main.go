package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"trackcoro/controller"
	"trackcoro/database"
	"trackcoro/objectstorage"
	"trackcoro/router"
)

var initialized = false
var ginLambda *ginadapter.GinLambda

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !initialized {
		database.ConnectToDB()
		defer database.DB.Close()
		database.MigrateSchema()
		objectstorage.InitializeS3Session()
		controller.InitializeControllers()
		r := router.InitializeRouter()
		ginLambda = ginadapter.New(r)
		initialized = true
	}
	return ginLambda.Proxy(req)
}

func main() {
	lambda.Start(Handler)
}
