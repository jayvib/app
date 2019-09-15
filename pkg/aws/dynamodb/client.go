package dynamodb

import (
	"fmt"
	sdk "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jayvib/clean-architecture/log"
	"github.com/jayvib/clean-architecture/pkg/aws/config"
	"github.com/jayvib/clean-architecture/pkg/aws/session"
	"sync"
)

const (
	servicename = "DynamoDB"
)

func New(conf config.Config) (*DynamoDB, error) {
	s, err := session.New(conf)
	if err != nil {
		return nil, err
	}
	svc := sdk.New(s)
	return &DynamoDB{
		client:      svc,
		logger:      log.NewStandardOutLogger(),
		prefix:      conf.DefaultPrefix,
		tables:      make(map[string]*Table),
		writeTables: make(map[string]struct{}),
	}, nil
}

type DynamoDB struct {
	client *sdk.DynamoDB

	logger      log.Logger
	prefix      string
	tablesMu    sync.RWMutex
	tables      map[string]*Table
	writeTables map[string]struct{}
}

func (svc *DynamoDB) CreateTable(blueprint *TableBlueprint) error {
	if blueprint.HashKey == nil {
		err := fmt.Errorf(".CreateTable: empty hash key in TableBlueprint")
		svc.Errorf("%s", err.Error())
		return err
	}

	origName := blueprint.name
	blueprint.name = svc.prefix + blueprint.name

	in := blueprint.CreateTableInput()
	out, err := svc.client.CreateTable(in)
	if err != nil {
		svc.Errorf("error on creating table: table=%s error=%s",
			blueprint.GetTableName(), err)
		blueprint.name = origName
		return err
	}

	blueprint = newTableBlueprintFromDescription(out.TableDescription)
	return nil
}

func (svc *DynamoDB) Errorf(fmt string, args ...interface{}) {
	svc.logger.Errorf(fmt, args...)
}
