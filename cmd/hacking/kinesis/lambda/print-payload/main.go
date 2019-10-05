package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)

type Person struct {
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
}

func Handle(ctx context.Context, p Person) (string, error){
	return fmt.Sprintf("Hello %s %s!", p.FirstName, p.LastName), nil
}

func main() {
	lambda.Start(Handle)
}
