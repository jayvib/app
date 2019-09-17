package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	authordynamo "github.com/jayvib/app/author/repository/dynamo"
	"github.com/jayvib/app/config"
	"github.com/jayvib/app/user"
	"github.com/jayvib/app/user/delivery/http"
	userdynamo "github.com/jayvib/app/user/repository/dynamo"
	"github.com/jayvib/app/user/usecase"
	"log"
)

// Reference Tutorial:
// https://www.alexedwards.net/blog/serverless-api-with-go-and-aws-lambda

const (
	region = "us-east-1"
)

type HandlerFunc func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func newGinLambdaAdapter(usecase user.Usecase) *ginadapter.GinLambda {
	r := newEngine(usecase)
	ginLambda := ginadapter.New(r)
	return ginLambda
}

func newDynamoDb() *dynamodb.DynamoDB {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	return dynamodb.New(sess, aws.NewConfig().WithRegion(region))
}

func newEngine(usecase user.Usecase) *gin.Engine {
	e := gin.Default()
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	http.RegisterHandlers(conf, e, usecase)
	return e
}

func newHandler(usecase user.Usecase) HandlerFunc {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		ginLambda := newGinLambdaAdapter(usecase)
		return ginLambda.ProxyWithContext(ctx, req)
	}
}

func main() {
	svc := newDynamoDb()
	userrepo := userdynamo.New(svc)
	authorrepo := authordynamo.New(svc)
	uc := usecase.New(userrepo, authorrepo)
	lambda.Start(newHandler(uc))
}
