// +build integration,elasticsearch

package http_test

import (
	estestutil "github.com/jayvib/app/internal/elasticsearch/testutil"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ElasticsearchSuite struct {
	suite.Suite
	url             string
	index           string
	mappingFilename string
	client          *elastic.Client
	migrate         *ElasticsearchMigration
}

func (e *ElasticsearchSuite) SetupSuite() {
	client, err := estestutil.Setup(e.url, e.index, e.mappingFilename)
	e.Require().NoError(err)
	e.client = client
}

// TODO: Implement migration
type ElasticsearchMigration struct {
	dataProviderFunc func() map[string]interface{}
	client           *elastic.Client
	index            string
}

func (m *ElasticsearchMigration) Up(t *testing.T) {
	estestutil.LoadSampleDataFromProvider(t, m.client, m.index, m.dataProviderFunc)
}

func (m *ElasticsearchMigration) Down(t *testing.T) {
}
