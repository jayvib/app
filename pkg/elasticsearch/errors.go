package elasticsearch

import "github.com/pkg/errors"

var (
	// FormatErr is the error used when a reader provide a non-json format content body.
	FormatErr         = errors.New("elasticsearch: format error should be in json")
	EmptyIndexNameErr = errors.New("elasticsearch: empty index name")
	NilReaderErr      = errors.New("elasticsearch: nil reader")
)
