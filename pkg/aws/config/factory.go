package config

import "github.com/aws/aws-sdk-go/aws"

// Use Factory
func New(ctype Type) *aws.Config {
	var conf *aws.Config
	b := NewDirector(nil)
	switch ctype {
	case ProductionConfig:
	case StagingConfig:
	default:
		configBuilder := newTestConfigBuilder()
		b.SetBuilder(configBuilder)
		b.Build()
		conf = configBuilder.GetConfig()
	}
	return conf
}
