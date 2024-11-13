package msp_tenant_users

import (
	"context"
	"fmt"
	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/msp/users"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"sort"
)

func NewMspManagedTenantUsersResource() resource.Resource { return &MspManagedTenantUsersResource{} }

type MspManagedTenantUsersResource struct {
	client *cdoClient.Client
}

func sortUsersToOrderInPlanData(users []User, planData *MspManagedTenantUsersResourceModel) *[]User {
	userOrder := make(map[string]int)
	for i, user := range planData.Users {
		userOrder[user.Username.ValueString()] = i
	}

	sort.Slice(users, func(i, j int) bool {
		return userOrder[users[i].Username.ValueString()] < userOrder[users[j].Username.ValueString()]
	})

	return &users
}

func (resource *MspManagedTenantUsersResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
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
						"id": schema.StringAttribute{
							MarkdownDescription: "Universally unique identifier of the user",
							Computed:            true,
						},
						"username": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The name of the user in CDO. This must be a valid e-mail address if the user is not an API-only user.",
						},
						"roles": schema.ListAttribute{
							Required:            true,
							MarkdownDescription: "The roles to assign to the user in the CDO tenant. Note: this list can only contain one entry.",
							ElementType:         types.StringType,
							Validators: []validator.List{
								listvalidator.ValueStringsAre(
									stringvalidator.OneOf("ROLE_READ_ONLY", "ROLE_ADMIN", "ROLE_SUPER_ADMIN", "ROLE_DEPLOY_ONLY", "ROLE_EDIT_ONLY", "ROLE_VPN_SESSIONS_MANAGER"),
								),
								listvalidator.SizeAtMost(1),
							},
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

func (resource *MspManagedTenantUsersResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Debug(ctx, "Adding users to the MSSP-managed CDO tenant")
	var planData MspManagedTenantUsersResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &planData)...)

	if response.Diagnostics.HasError() {
		return
	}

	createdUserDetails, err := resource.client.CreateUsersInMspManagedTenant(ctx, *resource.buildMspUsersInput(&planData))

	if err != nil {
		response.Diagnostics.AddError("failed to create users in MSP-managed tenant", err.Error())
		return
	}

	planData.Users = *sortUsersToOrderInPlanData(*resource.transformApiResponseToPlan(createdUserDetails), &planData)

	response.Diagnostics.Append(response.State.Set(ctx, &planData)...)
}

func (resource *MspManagedTenantUsersResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading users from MSP-managed CDO tenant")
	var stateData MspManagedTenantUsersResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &stateData)...)

	userDetails, err := resource.client.ReadUsersInMspManagedTenant(ctx, *resource.buildMspUsersInput(&stateData))
	if err != nil {
		response.Diagnostics.AddError("failed to read users in MSP-managed tenant", err.Error())
		return
	}

	stateData.Users = *sortUsersToOrderInPlanData(*resource.transformApiResponseToPlan(userDetails), &stateData)
	response.Diagnostics.Append(response.State.Set(ctx, &stateData)...)
}

func (resource *MspManagedTenantUsersResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (resource *MspManagedTenantUsersResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting users from MSP-managed CDO tenant")
	var stateData MspManagedTenantUsersResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &stateData)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := resource.deleteAllUsersInState(ctx, &stateData)
	if err != nil {
		response.Diagnostics.AddError("failed to delete users", err.Error())
	}
	stateData.Users = []User{}
	response.Diagnostics.Append(response.State.Set(ctx, &stateData)...)
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

func (resource *MspManagedTenantUsersResource) deleteAllUsersInState(ctx context.Context, stateData *MspManagedTenantUsersResourceModel) (interface{}, error) {
	var usernames []string
	for _, user := range stateData.Users {
		usernames = append(usernames, user.Username.ValueString())
	}
	deleteInput := users.MspDeleteUsersInput{
		TenantUid: stateData.TenantUid.ValueString(),
		Usernames: usernames,
	}
	return resource.client.DeleteUsersInMspManagedTenant(ctx, deleteInput)
}

func (resource *MspManagedTenantUsersResource) buildMspUsersInput(planData *MspManagedTenantUsersResourceModel) *users.MspUsersInput {
	var nativeUsers []users.UserDetails

	// 2. use plan data to create user and fill up rest of the model
	for _, user := range planData.Users {
		username := user.Username.ValueString()
		// Convert user.Roles to a slice of strings
		var roles []string
		for _, roleValue := range user.Roles.Elements() {
			if roleStr, ok := roleValue.(types.String); ok {
				roles = append(roles, roleStr.ValueString())
			}
		}
		apiOnlyUser := user.ApiOnlyUser.ValueBool()
		nativeUsers = append(nativeUsers, users.UserDetails{
			Username:    username,
			Roles:       roles,
			ApiOnlyUser: apiOnlyUser,
		})
	}

	return &users.MspUsersInput{
		TenantUid: planData.TenantUid.ValueString(),
		Users:     nativeUsers,
	}
}

func (resource *MspManagedTenantUsersResource) transformApiResponseToPlan(createdUserDetails *[]users.UserDetails) *[]User {
	var users []User
	for _, userDetails := range *createdUserDetails {
		var roles []attr.Value
		for _, role := range userDetails.Roles {
			roles = append(roles, types.StringValue(role))
		}
		users = append(users, User{
			Id:          types.StringValue(userDetails.Uid),
			Username:    types.StringValue(userDetails.Username),
			ApiOnlyUser: types.BoolValue(userDetails.ApiOnlyUser),
			Roles:       types.ListValueMust(types.StringType, roles),
		})
	}

	return &users
}
