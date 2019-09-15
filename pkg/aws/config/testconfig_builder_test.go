// +build unit

package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTestingConfigBuilder_GetConfig(t *testing.T) {
	testConfigBuilder := newTestConfigBuilder()
	d := NewDirector(testConfigBuilder)
	d.Build()
	conf := testConfigBuilder.GetConfig()
	assert.Equal(t, defaultEndpoint, aws.StringValue(conf.Endpoint))
	assert.Nil(t, conf.S3ForcePathStyle)
}
