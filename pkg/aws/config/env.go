package config

import "os"

const (
	envKeyRegion   = "AWS_REGION"
	envKeyEndpoint = "AWS_ENDPOINT"
)

func envRegion() string {
	return os.Getenv(envKeyRegion)
}

func envEndpoint() string {
	return os.Getenv(envKeyEndpoint)
}
