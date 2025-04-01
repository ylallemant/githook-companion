package globals

import (
	"net"
	"net/http"
	"time"
)

var DefaultApiClient = &http.Client{
	Timeout: time.Second * 5,
	Transport: &http.Transport{
		MaxIdleConns:          100,
		MaxConnsPerHost:       20,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
		IdleConnTimeout:       15 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		DisableCompression:    false,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 10 * time.Second,
		}).DialContext,
	},
}
