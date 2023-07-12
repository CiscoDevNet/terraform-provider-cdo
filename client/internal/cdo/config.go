package cdo

import (
	"log"
	"net/http"
	"time"
)

type Config struct {
	BaseUrl  string
	ApiToken string

	// these parameters apply to each request
	Retries int
	Delay   time.Duration
	Timeout time.Duration
}

const (
	DefaultRetries = 3
	DefaultDelay   = 3 * time.Second
	DefaultTimeout = 3 * time.Minute
)

var (
	DefaultLogger     = log.Default()
	DefaultHttpClient = http.DefaultClient
)

func NewConfigWithDefault(baseUrl string, apiToken string) *Config {
	return NewConfig(baseUrl, apiToken, DefaultRetries, DefaultDelay, DefaultTimeout)
}

func NewConfig(baseUrl string, apiToken string, retries int, delay time.Duration, timeout time.Duration) *Config {
	return &Config{
		BaseUrl:  baseUrl,
		ApiToken: apiToken,
		Retries:  retries,
		Delay:    delay,
		Timeout:  timeout,
	}
}
