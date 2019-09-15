package elasticsearch

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"io"
	"io/ioutil"
)

func CreateIndex(name string, r io.Reader) (result *elastic.IndicesCreateResult, err error) {
	if name == "" {
		return nil, EmptyIndexNameErr
	}
	if r == nil {
		return nil, NilReaderErr
	}
	mappingByte, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Is the content body a JSON format?
	err = isJSONFormat(mappingByte)
	if err != nil {
		return nil, err
	}

	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	return client.CreateIndex(name).BodyString(string(mappingByte)).Do(context.Background())
}

func isJSONFormat(bite []byte) error {
	var jsonStr map[string]interface{}
	err := json.Unmarshal(bite, &jsonStr)
	if err != nil {
		return err
	}
	return nil
}
