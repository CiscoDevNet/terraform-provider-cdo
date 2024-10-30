package msp_tenant_users

import (
	"context"
	"fmt"
	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/msp/users"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func NewMspManagedTenantUsersResource() resource.Resource { return &MspManagedTenantUsersResource{} }

type MspManagedTenantUsersResource struct {
	client *cdoClient.Client
}

func (r *MspManagedTenantUsersResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "Provides a resource to add users to an MSP managed tenant.",
		Attributes: map[string]schema.Attribute{
			"tenant_uid": schema.StringAttribute{
				MarkdownDescription: "Universally unique identifier of the tenant to which the users should be added.",
				Required:            true,
			},
			"users": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"username": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The name of the user in CDO. This must be a valid e-mail address if the user is not an API-only user.",
						},
						"role": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The role to assign to the user in the CDO tenant.",
						},
						"api_only_user": schema.BoolAttribute{
							Required:            true,
							MarkdownDescription: "Whether the user is an API-only user",
						},
					},
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.RequiresReplace(), // TODO stop destroying and re-adding once read endpoint added
					},
				},
				MarkdownDescription: "The list of users to be added to the tenant. You can add a maximum of 50 users at a time.",
				Required:            true,
			},
		},
	}
}

func (r *MspManagedTenantUsersResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Debug(ctx, "Adding users to the MSSP-managed CDO tenant")
	var planData MspManagedTenantUsersResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &planData)...)

	if response.Diagnostics.HasError() {
		return
	}

	_, err := r.createAllUsersInPlan(ctx, &planData)
	if err != nil {
		response.Diagnostics.AddError("failed to create users in MSP-managed tenant", err.Error())
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, &planData)...)
}

func (r *MspManagedTenantUsersResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading users from MSP-managed CDO tenant is a NOOP")
}

func (r *MspManagedTenantUsersResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (r *MspManagedTenantUsersResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting users from MSP-managed CDO tenant")
	var stateData MspManagedTenantUsersResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &stateData)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := r.deleteAllUsersInState(ctx, &stateData)
	if err != nil {
		response.Diagnostics.AddError("failed to delete users", err.Error())
	}
}

func (*MspManagedTenantUsersResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_msp_managed_tenant_users"
}

func (resource *MspManagedTenantUsersResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

	resource.client = client
}

func (r *MspManagedTenantUsersResource) deleteAllUsersInState(ctx context.Context, stateData *MspManagedTenantUsersResourceModel) (interface{}, error) {
	var usernames []string
	for _, user := range stateData.Users {
		usernames = append(usernames, user.Username.ValueString())
	}
	deleteInput := users.MspDeleteUsersInput{
		TenantUid: stateData.TenantUid.ValueString(),
		Usernames: usernames,
	}
	return r.client.DeleteUsersInMspManagedTenant(ctx, deleteInput)
}

func (r *MspManagedTenantUsersResource) createAllUsersInPlan(ctx context.Context, planData *MspManagedTenantUsersResourceModel) (*[]users.UserInput, *users.CreateError) {
	var nativeUsers []users.UserInput

	// 2. use plan data to create tenant and fill up rest of the model
	for _, user := range planData.Users {
		username := user.Username.ValueString()
		role := user.Role.ValueString()
		apiOnlyUser := user.ApiOnlyUser.ValueBool()
		nativeUsers = append(nativeUsers, users.UserInput{
			Username:    username,
			Role:        role,
			ApiOnlyUser: apiOnlyUser,
		})
	}

	// TODO we need endpoint to read users in an MSP-managed tenant
	return r.client.CreateUsersInMspManagedTenant(ctx, users.MspCreateUsersInput{
		TenantUid: planData.TenantUid.ValueString(),
		Users:     nativeUsers,
	})
}
