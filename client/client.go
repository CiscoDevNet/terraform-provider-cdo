// Package client provides API entrypoint, defines operations for the user.
// It simply forward requests and do nothing else.
package client

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/msp/tenants"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/msp/users"
	"net/http"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/connectoronboarding"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/sec"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/sec/seconboarding"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/duoadminpanel"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/settings"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/settings/tenantsettings"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa/asaconfig"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudfmc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd/cloudftdonboarding"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/genericssh"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/tenant"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/user"

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

func (c *Client) ReadDeviceByName(ctx context.Context, inp device.ReadByNameAndTypeInput) (*device.ReadOutput, error) {
	return device.ReadByNameAndType(ctx, c.client, inp)
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

func (c *Client) ReadCloudFtdByUid(ctx context.Context, inp cloudftd.ReadByUidInput) (*cloudftd.ReadOutput, error) {
	return cloudftd.ReadByUid(ctx, c.client, inp)
}

func (c *Client) ReadCloudFtdByName(ctx context.Context, inp cloudftd.ReadByNameInput) (*cloudftd.ReadOutput, error) {
	return cloudftd.ReadByName(ctx, c.client, inp)
}

func (c *Client) CreateCloudFtd(ctx context.Context, inp cloudftd.CreateInput) (*cloudftd.CreateOutput, error) {
	return cloudftd.Create(ctx, c.client, inp)
}

func (c *Client) UpdateCloudFtd(ctx context.Context, inp cloudftd.UpdateInput) (*cloudftd.UpdateOutput, error) {
	return cloudftd.Update(ctx, c.client, inp)
}

func (c *Client) DeleteCloudFtd(ctx context.Context, inp cloudftd.DeleteInput) (*cloudftd.DeleteOutput, error) {
	return cloudftd.Delete(ctx, c.client, inp)
}

func (c *Client) ReadUserByUsername(ctx context.Context, inp user.ReadByUsernameInput) (*user.ReadUserOutput, error) {
	return user.ReadByUsername(ctx, c.client, inp)
}

func (c *Client) ReadUserByUid(ctx context.Context, inp user.ReadByUidInput) (*user.ReadUserOutput, error) {
	return user.ReadByUid(ctx, c.client, inp)
}

func (c *Client) CreateUser(ctx context.Context, inp user.CreateUserInput) (*user.CreateUserOutput, error) {
	return user.Create(ctx, c.client, inp)
}

func (c *Client) DeleteUser(ctx context.Context, inp user.DeleteUserInput) (*user.DeleteUserOutput, error) {
	return user.Delete(ctx, c.client, inp)
}

func (c *Client) UpdateUser(ctx context.Context, inp user.UpdateUserInput) (*user.UpdateUserOutput, error) {
	return user.Update(ctx, c.client, inp)
}

func (c *Client) GenerateApiToken(ctx context.Context, inp user.GenerateApiTokenInput) (*user.ApiTokenResponse, error) {
	return user.GenerateApiToken(ctx, c.client, inp)
}

func (c *Client) RevokeApiToken(ctx context.Context, inp user.RevokeApiTokenInput) (*user.RevokeApiTokenOutput, error) {
	return user.RevokeApiToken(ctx, c.client, inp)
}

func (c *Client) CreateFtdOnboarding(ctx context.Context, inp cloudftdonboarding.CreateInput) (*cloudftdonboarding.CreateOutput, error) {
	return cloudftdonboarding.Create(ctx, c.client, inp)
}

func (c *Client) UpdateFtdOnboarding(ctx context.Context, inp cloudftdonboarding.UpdateInput) (*cloudftdonboarding.UpdateOutput, error) {
	return cloudftdonboarding.Update(ctx, c.client, inp)
}

func (c *Client) ReadFtdOnboarding(ctx context.Context, inp cloudftdonboarding.ReadInput) (*cloudftdonboarding.ReadOutput, error) {
	return cloudftdonboarding.Read(ctx, c.client, inp)
}

func (c *Client) DeleteFtdOnboarding(ctx context.Context, inp cloudftdonboarding.DeleteInput) (*cloudftdonboarding.DeleteOutput, error) {
	return cloudftdonboarding.Delete(ctx, c.client, inp)
}

func (c *Client) ReadTenantDetails(ctx context.Context) (*tenant.ReadTenantDetailsOutput, error) {
	return tenant.ReadTenantDetails(ctx, c.client)
}

func (c *Client) CreateCloudFmcDevice(ctx context.Context, inp cloudfmc.CreateInput) (*cloudfmc.CreateOutput, error) {
	return cloudfmc.Create(ctx, c.client, inp)
}

func (c *Client) ReadCloudFmcDevice(ctx context.Context) (*cloudfmc.ReadOutput, error) {
	return cloudfmc.Read(ctx, c.client, cloudfmc.NewReadInput())
}

func (c *Client) ReadCloudFmcSpecificDevice(ctx context.Context, inp cloudfmc.ReadSpecificInput) (*cloudfmc.ReadSpecificOutput, error) {
	return cloudfmc.ReadSpecific(ctx, c.client, inp)
}

func (c *Client) CreateConnectorOnboarding(ctx context.Context, inp connectoronboarding.CreateInput) (*connectoronboarding.CreateOutput, error) {
	return connectoronboarding.Create(ctx, c.client, inp)
}

func (c *Client) UpdateConnectorOnboarding(ctx context.Context, inp connectoronboarding.UpdateInput) (*connectoronboarding.UpdateOutput, error) {
	return connectoronboarding.Update(ctx, c.client, inp)
}

func (c *Client) ReadConnectorOnboarding(ctx context.Context, inp connectoronboarding.ReadInput) (*connectoronboarding.ReadOutput, error) {
	return connectoronboarding.Read(ctx, c.client, inp)
}

func (c *Client) DeleteConnectorOnboarding(ctx context.Context, inp connectoronboarding.DeleteInput) (*connectoronboarding.DeleteOutput, error) {
	return connectoronboarding.Delete(ctx, c.client, inp)
}

func (c *Client) CreateSec(ctx context.Context, inp sec.CreateInput) (*sec.CreateOutput, error) {
	return sec.Create(ctx, c.client, inp)
}

func (c *Client) UpdateSec(ctx context.Context, inp sec.UpdateInput) (*sec.UpdateOutput, error) {
	return sec.Update(ctx, c.client, inp)
}

func (c *Client) DeleteSec(ctx context.Context, inp sec.DeleteInput) (*sec.DeleteOutput, error) {
	return sec.Delete(ctx, c.client, inp)
}

func (c *Client) ReadSec(ctx context.Context, inp sec.ReadInput) (*sec.ReadOutput, error) {
	return sec.Read(ctx, c.client, inp)
}

func (c *Client) CreateSecOnboarding(ctx context.Context, inp seconboarding.CreateInput) (*seconboarding.CreateOutput, error) {
	return seconboarding.Create(ctx, c.client, inp)
}

func (c *Client) CreateDuoAdminPanel(ctx context.Context, inp duoadminpanel.CreateInput) (*duoadminpanel.CreateOutput, error) {
	return duoadminpanel.Create(ctx, c.client, inp)
}

func (c *Client) UpdateDuoAdminPanel(ctx context.Context, inp duoadminpanel.UpdateInput) (*duoadminpanel.UpdateOutput, error) {
	return duoadminpanel.Update(ctx, c.client, inp)
}

func (c *Client) ReadDuoAdminPanel(ctx context.Context, inp duoadminpanel.ReadByUidInput) (*duoadminpanel.ReadOutput, error) {
	return duoadminpanel.ReadByUid(ctx, c.client, inp)
}

func (c *Client) DeleteDuoAdminPanel(ctx context.Context, inp duoadminpanel.DeleteInput) (*duoadminpanel.DeleteOutput, error) {
	return duoadminpanel.Delete(ctx, c.client, inp)
}

func (c *Client) ReadTenantSettings(ctx context.Context) (*settings.TenantSettings, error) {
	return tenantsettings.Read(ctx, c.client)
}

func (c *Client) UpdateTenantSettings(ctx context.Context, updateTenantSettingsInput tenantsettings.UpdateTenantSettingsInput) (*settings.TenantSettings, error) {
	return tenantsettings.Update(ctx, c.client, updateTenantSettingsInput)
}

func (c *Client) CreateTenantUsingMspPortal(ctx context.Context, createInput tenants.MspCreateTenantInput) (*tenants.MspTenantOutput, *tenants.CreateError) {
	return tenants.Create(ctx, c.client, createInput)
}

func (c *Client) ReadMspManagedTenantByUid(ctx context.Context, readByUidInput tenants.ReadByUidInput) (*tenants.MspTenantOutput, error) {
	return tenants.ReadByUid(ctx, c.client, readByUidInput)
}

func (c *Client) DeleteMspManagedTenantByUid(ctx context.Context, deleteByUidInput tenants.DeleteByUidInput) (interface{}, error) {
	return tenants.DeleteByUid(ctx, c.client, deleteByUidInput)
}

func (c *Client) FindMspManagedTenantByName(ctx context.Context, readByNameInput tenants.ReadByNameInput) (*tenants.MspTenantsOutput, error) {
	return tenants.ReadByName(ctx, c.client, readByNameInput)
}

func (c *Client) CreateUsersInMspManagedTenant(ctx context.Context, createInput users.MspUsersInput) (*[]users.UserDetails, *users.CreateError) {
	return users.Create(ctx, c.client, createInput)
}

func (c *Client) ReadUsersInMspManagedTenant(ctx context.Context, readInput users.MspUsersInput) (*[]users.UserDetails, error) {
	return users.ReadCreatedUsersInTenant(ctx, c.client, readInput)
}

func (c *Client) DeleteUsersInMspManagedTenant(ctx context.Context, deleteInput users.MspDeleteUsersInput) (interface{}, error) {
	return users.Delete(ctx, c.client, deleteInput)
}
