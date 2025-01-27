package httpclient

import (
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var httpTransport *http.Transport
var otelHttpTransport *otelhttp.Transport

func CreateHttpClient() *http.Client {
	if httpTransport == nil {
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

func CreateOtelHttpClient() *http.Client {
	// client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	if otelHttpTransport == nil {
		transport := &http.Transport{
			MaxConnsPerHost:       10,
			MaxIdleConns:          10,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       15 * time.Second,
			ResponseHeaderTimeout: 15 * time.Second,
			DisableKeepAlives:     false,
		}
		otelHttpTransport = otelhttp.NewTransport(transport)
	}

	httpClient := &http.Client{
		Transport: otelHttpTransport,
		Timeout:   30 * time.Second,
	}

	return httpClient
}
