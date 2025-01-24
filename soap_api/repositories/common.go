package repositories

import (
	"net/http"
	"time"
)

var httpTransport *http.Transport

func CreateHttpClient() *http.Client {
	if httpTransport == nil {
		// transport = &http.Transport{
		// 	MaxIdleConns:          10,
		// 	IdleConnTimeout:       15 * time.Second,
		// 	ResponseHeaderTimeout: 15 * time.Second,
		// 	DisableKeepAlives:     false,
		// }
		t := http.DefaultTransport.(*http.Transport).Clone()
		t.MaxConnsPerHost = 10
		t.MaxIdleConns = 10
		t.MaxIdleConnsPerHost = 10
		t.IdleConnTimeout = 30 * time.Second
		t.ResponseHeaderTimeout = 30 * time.Second
		t.DisableKeepAlives = false
		// t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		httpTransport = t
	}

	httpClient := &http.Client{
		Transport: httpTransport,
		Timeout:   30 * time.Second,
	}

	return httpClient
}
