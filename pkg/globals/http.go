package globals

import (
	"net/http"
	"time"
)

var DefaultApiClient = &http.Client{
	Timeout: 5 * time.Second,
}
