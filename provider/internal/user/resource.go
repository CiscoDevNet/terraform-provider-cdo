package user

import (
	"context"
	"fmt"

	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/user"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	client *cdoClient.Client
}

type UserResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	GeneratedUsername types.String `tfsdk:"generated_username"`
	ApiOnlyUser       types.Bool   `tfsdk:"is_api_only_user"`
	UserRole          types.String `tfsdk:"role"`
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cdoClient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *cdoClient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a user resource. This allows a user to be created, updated, and deleted on CDO.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the user. This is a UUID and is automatically generated when the user is created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The username. If the user is not an [API only user](https://www.cisco.com/c/en/us/td/docs/security/cdo/managing-ftd-with-cdo/managing-ftd-with-cisco-defense-orchestrator/basics-of-cisco-defense-orchestrator.html?bookSearch=true#Cisco_Task.dita_d5ae397b-5aa5-4de0-82c1-a4aff63c5ba1), it must be an e-mail address; if the user is an API-only user, it must not be an email address, and CDO will generate a name for the user prefixed by the value provided here (see `generated_username`).",
				Required:            true,
			},
			"generated_username": schema.StringAttribute{
				MarkdownDescription: "The username generated by CDO. If the user is an [API only user](https://www.cisco.com/c/en/us/td/docs/security/cdo/managing-ftd-with-cdo/managing-ftd-with-cisco-defense-orchestrator/basics-of-cisco-defense-orchestrator.html?bookSearch=true#Cisco_Task.dita_d5ae397b-5aa5-4de0-82c1-a4aff63c5ba1), the username is appended with the name of the tenant (for example, an API-only user given the name `api_user` in the tenant `example` will have the generated username `api_user@CDO_example`). Otherwise, it is the same as the username entered in the `name` field.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_api_only_user": schema.BoolAttribute{
				MarkdownDescription: "Indicate whether this user is an [API only user](https://www.cisco.com/c/en/us/td/docs/security/cdo/managing-ftd-with-cdo/managing-ftd-with-cisco-defense-orchestrator/basics-of-cisco-defense-orchestrator.html?bookSearch=true#Cisco_Task.dita_d5ae397b-5aa5-4de0-82c1-a4aff63c5ba1)",
				Required:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "There are a variety of user roles in Cisco Defense Orchestrator (CDO). User roles are configured for each user on each tenant. See [User Roles in CDO](https://www.cisco.com/c/en/us/td/docs/security/cdo/managing-asa-with-cdo/managing-asa-with-cisco-defense-orchestrator/basics-of-cisco-defense-orchestrator.html#User_Roles) to learn more. Valid Values: (ROLE_READ_ONLY, ROLE_ADMIN, ROLE_SUPER_ADMIN, ROLE_DEPLOY_ONLY, ROLE_EDIT_ONLY, ROLE_VPN_SESSIONS_MANAGER)",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("ROLE_READ_ONLY", "ROLE_ADMIN", "ROLE_SUPER_ADMIN", "ROLE_DEPLOY_ONLY", "ROLE_EDIT_ONLY", "ROLE_VPN_SESSIONS_MANAGER"),
				},
			},
		},
	}
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	tflog.Trace(ctx, "create User resource")

	// 1. read terraform plan data into model
	var planData UserResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 2. create resource & fill model data
	createInp := user.NewCreateUserInput(planData.Name.ValueString(), planData.UserRole.ValueString(), planData.ApiOnlyUser.ValueBool())
	createUserOutp, err := r.client.CreateUser(ctx, *createInp)
	if err != nil {
		res.Diagnostics.AddError("failed to create user resource", err.Error())
		return
	}
	planData.ID = types.StringValue(createUserOutp.Uid)
	planData.GeneratedUsername = types.StringValue(createUserOutp.Name)

	// 3. fill terraform state using model data
	res.Diagnostics.Append(res.State.Set(ctx, &planData)...)
	tflog.Trace(ctx, "create user resource done")
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	tflog.Trace(ctx, "update user resource")

	// 1. read plan and state data from terraform
	var planData UserResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if res.Diagnostics.HasError() {
		return
	}
	var stateData UserResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if res.Diagnostics.HasError() {
		return
	}

	var userRoles []string
	userRoles = append(userRoles, planData.UserRole.ValueString())
	updateInput := user.UpdateUserInput{
		Uid:       planData.ID.ValueString(),
		UserRoles: userRoles,
	}
	// 2. update resource & state data
	userDetails, err := r.client.UpdateUser(ctx, updateInput)
	if err != nil {
		res.Diagnostics.AddError("failed to update user resource", err.Error())
		return
	}
	stateData.ID = types.StringValue(userDetails.Uid)
	stateData.Name = types.StringValue(userDetails.Name)
	stateData.UserRole = types.StringValue(userDetails.UserRoles[0])

	// 3. update terraform state with updated state data
	res.Diagnostics.Append(res.State.Set(ctx, &stateData)...)
	tflog.Trace(ctx, "update user resource done")
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	tflog.Trace(ctx, "delete user resource")

	// 1. read state data from terraform state
	var stateData UserResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 2. delete the resource
	deleteUserInput := user.DeleteUserInput{
		Uid: stateData.ID.ValueString(),
	}
	_, err := r.client.DeleteUser(ctx, deleteUserInput)
	if err != nil {
		res.Diagnostics.AddError("failed to delete User resource", err.Error())
	}
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// 1. read terraform plan data into the model
	var stateData UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 2. do read
	tflog.Debug(ctx, "Reading user: "+stateData.ID.ValueString())
	readOutp, err := r.client.ReadUserByUid(ctx, *user.NewReadByUidInput(stateData.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("failed to read user resource", err.Error())
		return
	}
	stateData.ID = types.StringValue(readOutp.Uid)
	stateData.GeneratedUsername = types.StringValue(readOutp.Name)
	stateData.UserRole = types.StringValue(readOutp.UserRoles[0]) // while our API technically allows multiple roles, our UI does not support multiple roles
	stateData.ApiOnlyUser = types.BoolValue(readOutp.ApiOnlyUser)

	// 3. save data into terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
	tflog.Trace(ctx, "read user resource done")
}
