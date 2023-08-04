// API entrypoint, defines operations for the user.
// It simply forward requests and do nothing else.
package client

import (
	"context"
	"net/http"

	"github.com/cisco-lockhart/go-client/connector/sdc"
	"github.com/cisco-lockhart/go-client/device"
	"github.com/cisco-lockhart/go-client/device/ios"

	"github.com/cisco-lockhart/go-client/device/asa"
	"github.com/cisco-lockhart/go-client/internal/device/asaconfig"
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

func (c *Client) ReadAllSdcs(ctx context.Context, inp sdc.ReadAllInput) (*sdc.ReadAllOutput, error) {
	return sdc.ReadAll(ctx, c.client, inp)
}

func (c *Client) ReadSdcByName(ctx context.Context, inp sdc.ReadByNameInput) (*sdc.ReadOutput, error) {
	return sdc.ReadByName(ctx, c.client, inp)
}

func (c *Client) ReadSdcByUid(ctx context.Context, inp sdc.ReadByUidInput) (*sdc.ReadOutput, error) {
	return sdc.ReadByUid(ctx, c.client, inp)
}

func (c *Client) ReadAsa(ctx context.Context, inp asa.ReadInput) (*asa.ReadOutput, error) {
	return asa.Read(ctx, c.client, inp)
}

func (c *Client) ReadDeviceByName(ctx context.Context, inp device.ReadByNameAndDeviceTypeInput) (*device.ReadOutput, error) {
	return device.ReadByNameAndDeviceType(ctx, c.client, inp)
}

func (c *Client) CreateAsa(ctx context.Context, inp asa.CreateInput) (*asa.CreateOutput, error) {
	return asa.Create(ctx, c.client, inp)
}

func (c *Client) UpdateAsa(ctx context.Context, inp asa.UpdateInput) (*asa.UpdateOutput, error) {
	return asa.Update(ctx, c.client, inp)
}

func (c *Client) DeleteAsa(ctx context.Context, inp asa.DeleteInput) (*asa.DeleteOutput, error) {
	return asa.Delete(ctx, c.client, inp)
}

func (c *Client) ReadIos(ctx context.Context, inp ios.ReadInput) (*ios.ReadOutput, error) {
	return ios.Read(ctx, c.client, inp)
}

func (c *Client) CreateIos(ctx context.Context, inp ios.CreateInput) (*ios.CreateOutput, error) {
	return ios.Create(ctx, c.client, inp)
}

func (c *Client) UpdateIos(ctx context.Context, inp ios.UpdateInput) (*ios.UpdateOutput, error) {
	return ios.Update(ctx, c.client, inp)
}

func (c *Client) DeleteIos(ctx context.Context, inp ios.DeleteInput) (*ios.DeleteOutput, error) {
	return ios.Delete(ctx, c.client, inp)
}

func (c *Client) ReadAsaConfig(ctx context.Context, inp asaconfig.ReadInput) (*asaconfig.ReadOutput, error) {
	return asaconfig.Read(ctx, c.client, inp)
}

func (c *Client) ReadSpecificAsa(ctx context.Context, inp asa.ReadSpecificInput) (*asa.ReadSpecificOutput, error) {
	return asa.ReadSpecific(ctx, c.client, inp)
}

func (c *Client) CreateSdc(ctx context.Context, inp sdc.CreateInput) (*sdc.CreateOutput, error) {
	return sdc.Create(ctx, c.client, inp)
}

func (c *Client) UpdateSdc(ctx context.Context, inp sdc.UpdateInput) (*sdc.UpdateOutput, error) {
	return sdc.Update(ctx, c.client, inp)
}

func (c *Client) DeleteSdc(ctx context.Context, inp sdc.DeleteInput) (*sdc.DeleteOutput, error) {
	return sdc.Delete(ctx, c.client, inp)
}
