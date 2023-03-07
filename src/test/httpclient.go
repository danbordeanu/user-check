package main

import (
	"net/http"
)

// funcHttpClientGet simple http get client
func funcHttpClientGet(urlStr string) (*http.Client, *http.Response) {
	// http client and request
	client := &http.Client{}
	r, _ := http.Get(urlStr)
	return client,r
}