package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func process(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	fileHeader := r.MultipartForm.File["uploaded"][0]
	file, _ := fileHeader.Open()
	data, _ := ioutil.ReadAll(file)
	fmt.Fprintln(w, string(data))
}

func main() {
	http.HandleFunc("/upload", process)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
