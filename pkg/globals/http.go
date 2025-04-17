package globals

import (
	"net"
	"net/http"
	"time"
)

var DefaultApiClient = &http.Client{
	// Transport: transport(),
}

func transport() *http.Transport {
	return &http.Transport{
		DialTLSContext: dialer().DialContext,
		DialContext:    dialer().DialContext,
	}
}

func dialer() *net.Dialer {
	return &net.Dialer{
		Timeout: 3 * time.Second,
	}
}
