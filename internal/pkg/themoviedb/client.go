package themoviedb

import (
	"net/http"
	"time"
)

type ClientConfig struct {
	ApiKey string
}

type ClientIns struct {
	httpClient *http.Client
	apiKey     string
}

func Client(config ClientConfig) ClientIns {
	clientIns := ClientIns{
		apiKey: config.ApiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}

	return clientIns
}
