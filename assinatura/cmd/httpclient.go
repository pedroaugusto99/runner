package cmd

import (
	"net/http"
	"time"
)

var signerHTTPClient = &http.Client{Timeout: 30 * time.Second}
