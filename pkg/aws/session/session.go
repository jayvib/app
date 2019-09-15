package session

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jayvib/clean-architecture/pkg/aws/config"
)

func New(conf config.Config) (*session.Session, error) {
	return session.NewSession(conf.AWSConfig())
}
