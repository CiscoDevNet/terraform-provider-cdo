package client

import "net/http"

// Client is used to communicate with the Cisco CDO Platform
type Client struct {
	httpClient *http.Client
	apiToken   string
}

// New instantiates a new Client with default HTTP configuration
func New(apiToken string) *Client {
	return NewWithHttpClient(&http.Client{}, apiToken)
}

// NewWithHttpClient instantiates a new Client with provided HTTP configuration
func NewWithHttpClient(httpClient *http.Client, apiToken string) *Client {
	return &Client{
		httpClient,
		apiToken,
	}
}
