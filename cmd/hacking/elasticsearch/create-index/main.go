package main

import (
	"io"
	"io/ioutil"
)

func readMappingTemplate(r io.Reader) (string, error) {
	byt, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(byt), nil
}

func main() {
	//file, err := os.Open("../../../user-mapping-es.json")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//mappingJson, err := readMappingTemplate(file)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//client, err := elastic.NewClient()
	//if err != nil {
	//	log.Fatal(err)
	//}
}
