// Package client provides API entrypoint, defines operations for the user.
// It simply forward requests and do nothing else.
package client

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa/asaconfig"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/genericssh"
	"net/http"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/ios"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa"
	internalhttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

type Client struct {
	client internalhttp.Client
}

// New instantiates a new Client with default HTTP configuration
func New(hostname, apiToken string) (*Client, error) {
	return NewWithHttpClient(http.DefaultClient, hostname, apiToken)
}

// NewWithHttpClient instantiates a new Client with provided HTTP configuration
func NewWithHttpClient(httpClient *http.Client, hostname, apiToken string) (*Client, error) {
	// log.SetOutput(os.Stdout)  // TODO: set this to os.Stdout in local environment
	client, err := internalhttp.NewWithHttpClient(httpClient, hostname, apiToken)
	if err != nil {
		return nil, err
	}
	return &Client{
		client: *client,
	}, nil
}

func (c *Client) ReadAllConnectors(ctx context.Context, inp connector.ReadAllInput) (*connector.ReadAllOutput, error) {
	return connector.ReadAll(ctx, c.client, inp)
}

func (c *Client) ReadConnectorByName(ctx context.Context, inp connector.ReadByNameInput) (*connector.ReadOutput, error) {
	return connector.ReadByName(ctx, c.client, inp)
}

func (c *Client) ReadConnectorByUid(ctx context.Context, inp connector.ReadByUidInput) (*connector.ReadOutput, error) {
	return connector.ReadByUid(ctx, c.client, inp)
}

func (c *Client) ReadAsa(ctx context.Context, inp asa.ReadInput) (*asa.ReadOutput, error) {
	return asa.Read(ctx, c.client, inp)
}

func (c *Client) ReadDeviceByName(ctx context.Context, inp device.ReadByNameAndDeviceTypeInput) (*device.ReadOutput, error) {
	return device.ReadByNameAndDeviceType(ctx, c.client, inp)
}

func (c *Client) CreateAsa(ctx context.Context, inp asa.CreateInput) (*asa.CreateOutput, *asa.CreateError) {
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

func (c *Client) CreateIos(ctx context.Context, inp ios.CreateInput) (*ios.CreateOutput, *ios.CreateError) {
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

func (c *Client) CreateConnector(ctx context.Context, inp connector.CreateInput) (*connector.CreateOutput, error) {
	return connector.Create(ctx, c.client, inp)
}

func (c *Client) UpdateConnector(ctx context.Context, inp connector.UpdateInput) (*connector.UpdateOutput, error) {
	return connector.Update(ctx, c.client, inp)
}

func (c *Client) DeleteConnector(ctx context.Context, inp connector.DeleteInput) (*connector.DeleteOutput, error) {
	return connector.Delete(ctx, c.client, inp)
}

func (c *Client) ReadGenericSSH(ctx context.Context, inp genericssh.ReadInput) (*genericssh.ReadOutput, error) {
	return genericssh.Read(ctx, c.client, inp)
}

func (c *Client) CreateGenericSSH(ctx context.Context, inp genericssh.CreateInput) (*genericssh.CreateOutput, error) {
	return genericssh.Create(ctx, c.client, inp)
}

func (c *Client) UpdateGenericSSH(ctx context.Context, inp genericssh.UpdateInput) (*genericssh.UpdateOutput, error) {
	return genericssh.Update(ctx, c.client, inp)
}

func (c *Client) DeleteGenericSSH(ctx context.Context, inp genericssh.DeleteInput) (*genericssh.DeleteOutput, error) {
	return genericssh.Delete(ctx, c.client, inp)
}
