package msp_tenant_user_groups

import (
	"context"
	"fmt"
	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/msp/usergroups"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"sort"
)

func NewMspManagedTenantUserGroupsResource() resource.Resource {
	return &MspManagedTenantUserGroupsResource{}
}

type MspManagedTenantUserGroupsResource struct {
	client *cdoClient.Client
}

func sortUserGroupsToOrderInPlanData(users []UserGroup, planData *MspManagedTenantUserGroupsResourceModel) *[]UserGroup {
	userOrder := make(map[string]int)
	for i, user := range planData.UserGroups {
		userOrder[user.GroupIdentifier.ValueString()] = i
	}

	sort.Slice(users, func(i, j int) bool {
		return userOrder[users[i].GroupIdentifier.ValueString()] < userOrder[users[j].GroupIdentifier.ValueString()]
	})

	return &users
}

func (resource *MspManagedTenantUserGroupsResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_msp_managed_tenant_user_groups"
}

func (resource *MspManagedTenantUserGroupsResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "Provides a resource to add user groups to an MSP managed tenant.",
		Attributes: map[string]schema.Attribute{
			"tenant_uid": schema.StringAttribute{
				MarkdownDescription: "Universally unique identifier of the tenant to which the user group should be added.",
				Required:            true,
			},
			"user_groups": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Universally unique identifier of the user group on Security Cloud Control",
							Computed:            true,
						},
						"group_identifier": schema.StringAttribute{
							MarkdownDescription: "The unique identifier of the user group in your Identity Provider (IdP)",
							Required:            true,
						},
						"issuer_url": schema.StringAttribute{
							MarkdownDescription: "The Identity Provider (IdP) URL, which Security Cloud Control will use to validate SAML assertions during the sign-in process",
							Required:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the user group. Security Cloud Control does not support special characters for this field.",
							Required:            true,
						},
						"notes": schema.StringAttribute{
							MarkdownDescription: "Any human-readable notes that are applicable to this user group",
							Optional:            true,
						},
						"role": schema.StringAttribute{
							MarkdownDescription: "The role assigned in Security Cloud Control to all users in this user group",
							Required:            true,
							Validators:          []validator.String{stringvalidator.OneOf("ROLE_READ_ONLY", "ROLE_ADMIN", "ROLE_SUPER_ADMIN", "ROLE_DEPLOY_ONLY", "ROLE_EDIT_ONLY", "ROLE_VPN_SESSIONS_MANAGER")},
						},
					},
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.RequiresReplace(),
					},
				},
				MarkdownDescription: "The list of user groups to be added to the tenant. You can add a maximum of 50 user groups at a time.",
				Required:            true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

func (resource *MspManagedTenantUserGroupsResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var planData MspManagedTenantUserGroupsResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &planData)...)
	if response.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to parse plan data")
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Adding user-group %v to MSSP-managed CDO tenant", planData))
	createdUserGroups, err := resource.client.CreateUserGroupsInMspManagedTenant(ctx, planData.TenantUid.ValueString(), resource.buildMspUserGroupInput(&planData))
	if err != nil {
		response.Diagnostics.AddError("Failed to create user group: %v", err.Error())
		return
	}

	planData.UserGroups = *sortUserGroupsToOrderInPlanData(*resource.transformApiResponseToPlan(createdUserGroups), &planData)

	response.Diagnostics.Append(response.State.Set(ctx, &planData)...)
}

func (resource *MspManagedTenantUserGroupsResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading user groups from MSP-managed CDO tenant")
	var stateData MspManagedTenantUserGroupsResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &stateData)...)
	userGroupDetails, err := resource.client.ReadUserGroupsInMspManagedTenant(ctx, stateData.TenantUid.ValueString(), resource.buildMspUserGroupInput(&stateData))
	if err != nil {
		response.Diagnostics.AddError("failed to read users in MSP-managed tenant", err.Error())
		return
	}
	stateData.UserGroups = *sortUserGroupsToOrderInPlanData(*resource.transformApiResponseToPlan(userGroupDetails), &stateData)
	response.Diagnostics.Append(response.State.Set(ctx, &stateData)...)
}

func (resource *MspManagedTenantUserGroupsResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (resource *MspManagedTenantUserGroupsResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Debug(ctx, "Deleting user groups from MSP-managed CDO tenant")
	var stateData MspManagedTenantUserGroupsResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &stateData)...)
	if response.Diagnostics.HasError() {
		return
	}
	_, err := resource.deleteAllUserGroupsInState(ctx, &stateData)
	if err != nil {
		response.Diagnostics.AddError("failed to delete users", err.Error())
	}

	stateData.UserGroups = []UserGroup{}
	response.Diagnostics.Append(response.State.Set(ctx, &stateData)...)
}

func (resource *MspManagedTenantUserGroupsResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

func (resource *MspManagedTenantUserGroupsResource) buildMspUserGroupInput(planData *MspManagedTenantUserGroupsResourceModel) *[]usergroups.MspManagedUserGroupInput {
	var userGroupCreateOrUpdateInput []usergroups.MspManagedUserGroupInput
	for _, userGroup := range planData.UserGroups {
		var notes *string
		if !userGroup.Notes.IsNull() {
			notes = userGroup.Notes.ValueStringPointer()
		}
		userGroupCreateOrUpdateInput = append(userGroupCreateOrUpdateInput, usergroups.MspManagedUserGroupInput{
			GroupIdentifier: userGroup.GroupIdentifier.ValueString(),
			IssuerUrl:       userGroup.IssuerUrl.ValueString(),
			Name:            userGroup.Name.ValueString(),
			Role:            userGroup.Role.ValueString(),
			Notes:           notes,
		})
	}

	return &userGroupCreateOrUpdateInput
}

func (resource *MspManagedTenantUserGroupsResource) transformApiResponseToPlan(createdUserGroupDetails *[]usergroups.MspManagedUserGroup) *[]UserGroup {
	var userGroups []UserGroup
	for _, userGroupDetails := range *createdUserGroupDetails {
		var notes basetypes.StringValue
		if userGroupDetails.Notes != nil {
			notes = types.StringValue(*userGroupDetails.Notes)
		}
		userGroups = append(userGroups, UserGroup{
			Id:              types.StringValue(userGroupDetails.Uid),
			GroupIdentifier: types.StringValue(userGroupDetails.GroupIdentifier),
			IssuerUrl:       types.StringValue(userGroupDetails.IssuerUrl),
			Name:            types.StringValue(userGroupDetails.Name),
			Role:            types.StringValue(userGroupDetails.Role),
			Notes:           notes,
		})
	}

	return &userGroups
}

func (resource *MspManagedTenantUserGroupsResource) deleteAllUserGroupsInState(ctx context.Context, stateData *MspManagedTenantUserGroupsResourceModel) (interface{}, error) {
	var userGroupUids []string
	for _, user := range stateData.UserGroups {
		userGroupUids = append(userGroupUids, user.Id.ValueString())
	}
	deleteInput := usergroups.MspManagedUserGroupDeleteInput{
		UserGroupUids: userGroupUids,
	}
	return resource.client.DeleteUserGroupsInMspManagedTenant(ctx, stateData.TenantUid.ValueString(), &deleteInput)
}
