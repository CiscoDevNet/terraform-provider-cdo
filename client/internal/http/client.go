// Client wrap Request and Respose to provide a slightly higher level API for internal use
package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/cisco-lockhart/go-client/internal/cdo"
)

type Client struct {
	config     cdo.Config
	httpClient *http.Client
	Logger     *log.Logger
}

// NewWithDefault instantiates a new Client with default HTTP configuration
func NewWithDefault(baseUrl, apiToken string) *Client {
	return NewWithHttpClient(cdo.DefaultHttpClient, baseUrl, apiToken)
}

// NewWithHttpClient instantiates a new Client with provided HTTP configuration
func NewWithHttpClient(client *http.Client, baseUrl, apiToken string) *Client {
	return New(client, cdo.DefaultLogger, baseUrl, apiToken, cdo.DefaultRetries, cdo.DefaultDelay, cdo.DefaultTimeout)
}

func New(
	// objects
	client *http.Client,
	logger *log.Logger,

	// configs
	baseUrl string,
	apiToken string,
	retries int,
	delay time.Duration,
	timeout time.Duration,

) *Client {
	return &Client{
		config:     *cdo.NewConfig(baseUrl, apiToken, retries, delay, timeout),
		httpClient: client,
		Logger:     logger,
	}
}

func (c *Client) NewGet(ctx context.Context, url string) *Request {
	return NewRequest(c.config, c.httpClient, c.Logger, ctx, "GET", url, nil)
}

func (c *Client) NewDelete(ctx context.Context, url string) *Request {
	return NewRequest(c.config, c.httpClient, c.Logger, ctx, "DELETE", url, nil)
}

func (c *Client) NewPost(ctx context.Context, url string, body any) *Request {
	return NewRequest(c.config, c.httpClient, c.Logger, ctx, "POST", url, body)
}

func (c *Client) NewPut(ctx context.Context, url string, body any) *Request {
	return NewRequest(c.config, c.httpClient, c.Logger, ctx, "PUT", url, body)
}

func (c *Client) BaseUrl() string {
	return c.config.BaseUrl
}
