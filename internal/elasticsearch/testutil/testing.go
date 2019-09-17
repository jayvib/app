package testutil

import (
	"context"
	"fmt"
	"github.com/jayvib/app/log"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

const elasticsearchTestURL = "http://localhost:9200"

const (
	mappingDirectory = "testdata"
)

// SetupHelper use to do setup operation for elasticsearch integration test it wraps
// the Setup function internally.
func SetupHelper(t *testing.T, url, indexName, mappingFilename string) (client *elastic.Client) {
	t.Helper()
	client, err := Setup(url, indexName, mappingFilename)
	require.NoError(t, err)
	return client
}

// Setup is an helper that will be use for testing.
//
// It will create a new index indexName using the mapping provided through
// mappingFilename and returns the initiated elasticsearch client and an
// error if has any.
//
// The mapping file will be look under the testdata folder, so make sure
// that the file exist under the testdata else it will return a os.PathError
// error
func Setup(url string, indexName, mappingFilename string) (client *elastic.Client, err error) {
	if url == "" {
		url = elasticsearchTestURL
	}
	// Be transparent and explicit for establishing
	// a client connection.
	client, err = elastic.NewClient(
		elastic.SetURL(elasticsearchTestURL))
	if err != nil {
		return nil, err
	}

	err = createIndex(client, indexName, mappingFilename)
	if err != nil {
		return
	}
	return client, nil
}

func createIndex(client *elastic.Client, indexName, mappingFilename string) error {
	// Setup the index
	exist, err := client.IndexExists(indexName).Do(context.Background())
	if err != nil {
		return err
	}
	if !exist {
		mapping, err := ioutil.ReadFile(filepath.Join(mappingDirectory, mappingFilename))
		if err != nil {
			return err
		}
		log.Infof("mapping:", string(mapping))
		result, err := client.
			CreateIndex(indexName).
			BodyString(string(mapping)).
			Do(context.Background())
		if err != nil {
			return err
		}

		if !result.Acknowledged {
			return fmt.Errorf("creating user index not acknowledged\n")
		}
	}
	return nil
}

func LoadSampleDataFromProvider(t *testing.T, client *elastic.Client, index string, provider func() map[string]interface{}) {
	t.Helper()
	defer func() {
		_, err := client.Flush(index).Do(context.Background())
		require.NoError(t, err)
	}()
	data := provider()
	for id, body := range data {
		_, err := client.Index().Index(index).Id(id).BodyJson(body).Do(context.Background())
		require.NoError(t, err)
	}
	time.Sleep(2 * time.Second)
}
