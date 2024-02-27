package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/michgur/puncher/app"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambda

func Handler(context context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {
		ginLambda = ginadapter.New(app.R)
	}

	return ginLambda.ProxyWithContext(context, request)
}

func main() {
	if os.Getenv("ENV") == "lambda" {
		lambda.Start(Handler)
	} else {
		app.R.Run(":3000")
	}
}
