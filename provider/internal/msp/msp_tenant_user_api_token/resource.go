package msp_tenant_user_api_token

import (
	"context"
	"fmt"
	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/msp/users"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func NewMspManagedTenantUserApiTokenResource() resource.Resource {
	return &MspManagedTenantUserApiTokenResource{}
}

type MspManagedTenantUserApiTokenResource struct {
	client *cdoClient.Client
}

func (m *MspManagedTenantUserApiTokenResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_msp_managed_tenant_user_api_token"
}

func (m *MspManagedTenantUserApiTokenResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "Provides a resource to manage an API token for a user in an MSP-managed tenant.",
		Attributes: map[string]schema.Attribute{
			"tenant_uid": schema.StringAttribute{
				MarkdownDescription: "Universally unique identifier of the tenant in which the API token for the user should be generated.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_uid": schema.StringAttribute{
				MarkdownDescription: "Universally unique identifier of the user for whom the API token should be generated.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "The generated API token for the user.",
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (m *MspManagedTenantUserApiTokenResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var planData MspManagedTenantUserApiTokenResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &planData)...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Generating an API token for a user %s in the tenant %s...", planData.UserUid, planData.TenantUid))

	apiTokenInfo, err := m.client.GenerateApiTokenForUserInMspManagedTenant(ctx, users.MspGenerateApiTokenInput{
		UserUid:   planData.UserUid.ValueString(),
		TenantUid: planData.TenantUid.ValueString(),
	})
	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("Failed to generate API token for user %s in MSP-managed tenant %s", planData.UserUid, planData.TenantUid), err.Error())
		return
	}

	planData.ApiToken = types.StringValue(apiTokenInfo.ApiToken)
	response.Diagnostics.Append(response.State.Set(ctx, &planData)...)
}

func (m *MspManagedTenantUserApiTokenResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Debug(ctx, "This is a NOOP")
}

func (m *MspManagedTenantUserApiTokenResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (m *MspManagedTenantUserApiTokenResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var stateData MspManagedTenantUserApiTokenResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &stateData)...)
	if response.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Revoking an API token for a user %s in the tenant %s...", stateData.UserUid, stateData.TenantUid))

	_, err := m.client.RevokeApiTokenForUserInMspManagedTenant(ctx, users.MspRevokeApiTokenInput{ApiToken: stateData.ApiToken.ValueString()})
	if err != nil {
		response.Diagnostics.AddError("Failed to revoke API token for user", err.Error())
	}
}

func (m *MspManagedTenantUserApiTokenResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cdoClient.Client)

	if !ok {
		res.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *cdoClient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	m.client = client
}
