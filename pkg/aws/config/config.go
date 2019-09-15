package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

const (
	defaultRegion   = "us-east-1"
	defaultEndpoint = "http://localhost:8000"
)

// Config represents the aws configuration
// setting.
type Config struct {
	AccessKey        string
	SecretKey        string
	Region           string
	Endpoint         string
	Filename         string
	Profile          string
	DefaultPrefix    string
	S3ForcePathStyle bool
}

func (c Config) AWSConfig() *aws.Config {
	cred := c.awsCredentials()
	conf := &aws.Config{
		Credentials: cred,
		Region:      aws.String(c.getRegion()),
	}

	if ep := c.getEndpoint(); ep != "" {
		conf.Endpoint = aws.String(ep)
	}

	if c.S3ForcePathStyle {
		conf.S3ForcePathStyle = aws.Bool(true)
	}

	return conf
}

func (c Config) awsCredentials() *credentials.Credentials {
	// Get from environment
	cred := credentials.NewEnvCredentials()
	_, err := cred.Get()
	if err == nil {
		return cred
	}

	cred = credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, "")
	_, err = cred.Get()
	if err == nil {
		return cred
	}
	return credentials.NewSharedCredentials(c.Filename, c.Profile)
}

func (c Config) getRegion() string {
	// If the no region in the field use
	// the default region
	region := defaultRegion
	if c.Region != "" {
		region = c.Region
	}
	reg := envRegion()
	if reg != "" {
		return reg
	}
	return region
}

func (c Config) getEndpoint() string {
	if c.Endpoint != "" {
		return c.Endpoint
	}
	return envEndpoint()
}
