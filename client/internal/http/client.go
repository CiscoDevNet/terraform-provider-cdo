// Package http provides a Client that wrap Request and Response to provide a slightly higher level API for internal use
package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/cdo"
)

type Client struct {
	config     cdo.Config
	httpClient *http.Client
	Logger     *log.Logger
}

// NewWithDefault instantiates a new Client with default HTTP configuration
func NewWithDefault(baseUrl, apiToken string) (*Client, error) {
	return NewWithHttpClient(cdo.DefaultHttpClient, baseUrl, apiToken)
}

// NewWithHttpClient instantiates a new Client with provided HTTP configuration
func NewWithHttpClient(client *http.Client, baseUrl, apiToken string) (*Client, error) {
	return New(client, cdo.DefaultLogger, baseUrl, apiToken, cdo.DefaultRetries, cdo.DefaultDelay, cdo.DefaultTimeout)
}

func MustNewWithDefault(baseUrl, apiToken string) *Client {
	return MustNewWithHttpClient(cdo.DefaultHttpClient, baseUrl, apiToken)
}

func MustNewWithConfig(baseUrl, apiToken string, retries int, delay, timeout time.Duration) *Client {
	return MustNew(cdo.DefaultHttpClient, cdo.DefaultLogger, baseUrl, apiToken, retries, delay, timeout)
}

func MustNewWithHttpClient(client *http.Client, baseUrl, apiToken string) *Client {
	return MustNew(client, cdo.DefaultLogger, baseUrl, apiToken, cdo.DefaultRetries, cdo.DefaultDelay, cdo.DefaultTimeout)
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

) (*Client, error) {
	config, err := cdo.NewConfig(baseUrl, apiToken, retries, delay, timeout)
	if err != nil {
		return &Client{}, err
	}
	return &Client{
		config:     config,
		httpClient: client,
		Logger:     logger,
	}, nil
}

func MustNew(
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
	c, err := New(client, logger, baseUrl, apiToken, retries, delay, timeout)
	if err != nil {
		panic(fmt.Sprint("failed to create http client, cause=%w", err))
	}
	return c
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

func (c *Client) Host() string {
	return c.config.Host
}
