package main

// Use to build the lambda
//go:generate env GOOS=linux GOARCH=amd64 go build -o ./bin/main .
//go:generate zip -j ./bin/main.zip ./bin/main
