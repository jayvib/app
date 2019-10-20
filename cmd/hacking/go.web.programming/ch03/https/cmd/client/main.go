package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	pemfile, err := os.Open("cert.pem")
	if err != nil {
		log.Fatal(err)
	}

	pemBytes, _ := ioutil.ReadAll(pemfile)

	// create a certificate pool
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(pemBytes)
	clientTLSConf := &tls.Config{
		RootCAs: certPool,
	}

	transport := &http.Transport{
		TLSClientConfig: clientTLSConf,
	}

	client := http.Client{
		Transport: transport,
	}

	resp, err := client.Get("https://localhost:8080/")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Status:", resp.StatusCode)
	io.Copy(os.Stdout, resp.Body)

}
