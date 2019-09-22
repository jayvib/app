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
	usersearches "github.com/jayvib/app/user/search/elasticsearch"
	"github.com/jayvib/app/user/usecase"
	"github.com/olivere/elastic/v7"
	"log"
	"sync"
)

// Reference Tutorial:
// https://www.alexedwards.net/blog/serverless-api-with-go-and-aws-lambda

const (
	region = "us-east-1"
)

var (
	once sync.Once
	conf *config.Config
	db   *dynamodb.DynamoDB
	se   *elastic.Client
)

var (
	ginLambdaOnce sync.Once
	ginLambda     *ginadapter.GinLambda
)

func init() {
}

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
		ginLambdaOnce.Do(func() { // avoid re-initializing the gin adapter in the next invocation.
			ginLambda = newGinLambdaAdapter(usecase)
		})
		return ginLambda.ProxyWithContext(ctx, req)
	}
}

func newESClient() *elastic.Client {
	esClient, err := elastic.NewClient(
		elastic.SetURL(conf.Elasticsearch.Servers...),
		elastic.SetSniff(false))
	if err != nil {
		log.Fatal(err)
	}
	return esClient
}

func main() {
	var err error

	once.Do(func() {
		conf, err = config.New()
		if err != nil {
			log.Fatal(err)
		}
		db = newDynamoDb()
		se = newESClient()
	})
	userRepo := userdynamo.New(db)
	userSearchEngine := usersearches.New(se)
	authorrepo := authordynamo.New(db)

	// duh! Elasticsearch in lambda function??
	// TODO: Use kinesis to put the data that will be written to
	// elasticsearch.
	//
	// Elasticsearch will be running in EC2 machines under private network
	// so make sure that the lambda has the VPC-ID in order to access
	// the EC2 machines.
	uc := usecase.New(userRepo, authorrepo, userSearchEngine) // this will panic when run in the production
	lambda.Start(newHandler(uc))
}
