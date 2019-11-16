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

// Pool type is to reuse memory between garbage collections
var buffPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func MarshalJSON(p *Person) string{
	buff := buffPool.Get().(*bytes.Buffer)
	defer func() {
		buff.Reset()
		buffPool.Put(buff)
	}()
	err := json.NewEncoder(buff).Encode(p)
	if err != nil {
		log.Errorf("error: %v", err)
		return ""
	}
	return buff.String()
}
