package syncpool

import (
	"bytes"
	"encoding/json"
	"github.com/jayvib/app/log"
	"sync"
)

type Person struct {
	Name string `json:"name"`
}

var buffPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func MarshalJSON(p *Person) string{
	var buff bytes.Buffer
	err := json.NewEncoder(&buff).Encode(p)
	if err != nil {
		log.Errorf("error: %v", err)
		return ""
	}
	return buff.String()
}
