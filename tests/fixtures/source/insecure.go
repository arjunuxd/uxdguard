package main

import (
	"crypto/tls"
	"net/http"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	_ = tr
	client := &http.Client{}
	resp, _ := client.Get("http://example.com/api")
	_ = resp
}
