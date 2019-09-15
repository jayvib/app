package config

import "github.com/aws/aws-sdk-go/aws"

type Type int

const (
	LocalTestConfig = Type(iota)
	StagingConfig
	ProductionConfig
)

// A best explanation why use builder pattern.
// https://stackoverflow.com/questions/328496/when-would-you-use-the-builder-pattern
type Builder interface {
	SetAccessKey() Builder
	SetSecretKey() Builder
	SetRegion() Builder
	SetEndpoint() Builder
	SetFilename() Builder
	SetProfile() Builder
	SetDefaultPrefix() Builder
	SetS3ForcePathStyle() Builder
	GetConfig() *aws.Config
}

func NewDirector(b Builder) *Director {
	return &Director{
		builder: b,
	}
}

type Director struct {
	builder Builder
}

func (d *Director) Build() {
	d.builder.
		SetAccessKey().
		SetSecretKey().
		SetDefaultPrefix().
		SetEndpoint().
		SetFilename().
		SetProfile().
		SetRegion().
		SetS3ForcePathStyle()
}

func (d *Director) SetBuilder(b Builder) {
	d.builder = b
}
