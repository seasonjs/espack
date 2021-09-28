package http_cli

import (
	"net/http"
	"time"
)

var _httpCli = &http.Client{
	Timeout: time.Duration(15) * time.Second,
	Transport: &http.Transport{
		MaxIdleConnsPerHost:   1,
		MaxConnsPerHost:       1,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}
