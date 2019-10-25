package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/set_message", func(w http.ResponseWriter, r *http.Request) {
		msg := []byte("Hello World!")
		cookie := http.Cookie{
			Name:  "flash",
			Value: base64.StdEncoding.EncodeToString(msg),
		}

		http.SetCookie(w, &cookie)
	})

	http.HandleFunc("/get_message", func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("flash")
		rc := http.Cookie{
			Name:    "flash",
			MaxAge:  -1,
			Expires: time.Unix(1, 9),
		}

		http.SetCookie(w, &rc)
		val, _ := base64.StdEncoding.DecodeString(c.Value)
		fmt.Fprintln(w, string(val))
	})

	http.ListenAndServe(":8080", nil)
}
