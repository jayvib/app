package elasticsearch

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

// CreateIndex create a new elasticsearch index if not yet exists.
func CreateIndex(client *elastic.Client, indexName, mappingFilename string) error {
	// Setup the index
	exist, err := client.IndexExists(indexName).Do(context.Background())
	if err != nil {
		return err
	}
	if !exist {
		mapping, err := ioutil.ReadFile(mappingFilename)
		if err != nil {
			return err
		}
		logrus.Infof("Creating index: %s\n", indexName)
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

