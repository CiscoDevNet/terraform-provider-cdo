// API entrypoint, defines operations for the user.
// It simply forward requests and do nothing else.
package client

import (
	"context"
	"net/http"

	"github.com/cisco-lockhart/go-client/device/asa"
	"github.com/cisco-lockhart/go-client/device/asa/asaconfig"
	"github.com/cisco-lockhart/go-client/device/sdc"
	internalhttp "github.com/cisco-lockhart/go-client/internal/http"
)

type Client struct {
	client internalhttp.Client
}

// New instantiates a new Client with default HTTP configuration
func New(hostname, apiToken string) *Client {
	return NewWithHttpClient(http.DefaultClient, hostname, apiToken)
}

// NewWithHttpClient instantiates a new Client with provided HTTP configuration
func NewWithHttpClient(httpClient *http.Client, hostname, apiToken string) *Client {
	// log.SetOutput(os.Stdout)  // TODO: set this to os.Stdout in local environment
	return &Client{
		client: *internalhttp.NewWithHttpClient(httpClient, hostname, apiToken),
	}
}

func (c *Client) ReadAllSdcs(ctx context.Context, r sdc.ReadAllInput) (*sdc.ReadAllOutput, error) {
	return sdc.ReadAll(ctx, c.client, r)
}

func (c *Client) ReadAsa(ctx context.Context, r asa.ReadInput) (*asa.ReadOutput, error) {
	return asa.Read(ctx, c.client, r)
}

func (c *Client) CreateAsa(ctx context.Context, r asa.CreateInput) (*asa.CreateOutput, error) {
	return asa.Create(ctx, c.client, r)
}

func (c *Client) UpdateAsa(ctx context.Context, r asa.UpdateInput) (*asa.UpdateOutput, error) {
	return asa.Update(ctx, c.client, r)
}

func (c *Client) DeleteAsa(ctx context.Context, r asa.DeleteInput) (*asa.DeleteOutput, error) {
	return asa.Delete(ctx, c.client, r)
}

func (c *Client) ReadAsaConfig(ctx context.Context, r asaconfig.ReadInput) (*asaconfig.ReadOutput, error) {
	return asaconfig.Read(ctx, c.client, r)
}

func (c *Client) ReadSpecific(ctx context.Context, r asa.ReadSpecificInput) (*asa.ReadSpecificOutput, error) {
	return asa.ReadSpecific(ctx, c.client, r)
}
