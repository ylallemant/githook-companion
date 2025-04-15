package globals

import (
	"net"
	"net/http"
	"time"
)

var DefaultApiClient = &http.Client{
	Timeout:   time.Second * 5,
	Transport: transport(),
}

func transport() *http.Transport {
	return &http.Transport{
		MaxIdleConns:          100,
		MaxConnsPerHost:       20,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
		IdleConnTimeout:       5 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		DisableCompression:    false,
		DialTLSContext:        dialer().DialContext,
		DialContext:           dialer().DialContext,
	}
}

func dialer() *net.Dialer {
	return &net.Dialer{
		Timeout:   3 * time.Second,
		KeepAlive: 10 * time.Second,
	}
}
