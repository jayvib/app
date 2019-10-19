package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	file, err := os.Open("example.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	body := &bytes.Buffer{}
	mwriter := multipart.NewWriter(body)
	part, err := mwriter.CreateFormFile("file", "example.csv")
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(file)
	contentType := http.DetectContentType(content)
	fmt.Println(contentType)

	_, err = part.Write(content)
	if err != nil {
		log.Fatal(err)
	}
	err = mwriter.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(body.String())
}
