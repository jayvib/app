package session

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func New(conf aws.Config) *session.Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: conf,
	}))
	return sess
}
