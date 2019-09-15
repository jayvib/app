package config

import "github.com/aws/aws-sdk-go/aws"

func newTestConfigBuilder() *TestingConfigBuilder {
	return &TestingConfigBuilder{}
}

type TestingConfigBuilder struct {
	conf Config
}

func (tc *TestingConfigBuilder) SetAccessKey() Builder {
	tc.conf.AccessKey = "accesskey"
	return tc
}

func (tc *TestingConfigBuilder) SetSecretKey() Builder {
	tc.conf.SecretKey = "secretkey"
	return tc
}

func (tc *TestingConfigBuilder) SetRegion() Builder {
	tc.conf.Region = defaultRegion
	return tc
}

func (tc *TestingConfigBuilder) SetEndpoint() Builder {
	tc.conf.Endpoint = defaultEndpoint
	return tc
}

func (tc *TestingConfigBuilder) SetFilename() Builder {
	tc.conf.Filename = "config.test.json"
	return tc
}

func (tc *TestingConfigBuilder) SetProfile() Builder {
	tc.conf.Profile = "test-aws"
	return tc
}

func (tc *TestingConfigBuilder) SetDefaultPrefix() Builder {
	tc.conf.DefaultPrefix = "aws"
	return tc
}

func (tc *TestingConfigBuilder) SetS3ForcePathStyle() Builder {
	tc.conf.S3ForcePathStyle = false
	return tc
}

func (tc *TestingConfigBuilder) GetConfig() *aws.Config {
	return tc.conf.AWSConfig()
}
