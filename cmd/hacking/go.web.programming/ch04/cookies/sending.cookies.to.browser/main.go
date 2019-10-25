package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/set-cookie", func(w http.ResponseWriter, r *http.Request){
		c1 := http.Cookie{
			Name: "first_cookie",
			Value: "Go Web Programming",
			HttpOnly: true,
		}

		c2 := http.Cookie{
			Name: "second_cookie",
			Value: "Manning Publication Co",
			HttpOnly: true,
		}

		w.Header().Set("Set-Cookie", c1.String())
		w.Header().Set("Set-Cookie", c2.String())
	})

	http.HandleFunc("/get-cookie", func(w http.ResponseWriter, r *http.Request){
		cookie, _ := r.Cookie("second_cookie")
		fmt.Fprint(w, cookie)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

