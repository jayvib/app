package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// This code is an attempt to review host to
// read body from the request

func main() {
	http.HandleFunc("/body", func(w http.ResponseWriter, r *http.Request){
		len := r.ContentLength

		content, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		fmt.Println("Body:", string(content))
		fmt.Println("Len:", len)

	})
	http.ListenAndServe(":8080", nil)
}
