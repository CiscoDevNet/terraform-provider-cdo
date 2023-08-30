package user_api_token

import (
	"context"
	"fmt"

	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/user"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

type ApiTokenResourceModel struct {
	Username types.String `tfsdk:"username"`
	ApiToken types.String `tfsdk:"api_token"`
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
	resp.TypeName = req.ProviderTypeName + "_api_token"
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a resource to generate an API token for an API-only user. This allows an API-only user's token to be created and refreshed on CDO.",

		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				MarkdownDescription: "The username for which the API token should be generated. This username must be of an [API only user](https://www.cisco.com/c/en/us/td/docs/security/cdo/managing-ftd-with-cdo/managing-ftd-with-cisco-defense-orchestrator/basics-of-cisco-defense-orchestrator.html?bookSearch=true#Cisco_Task.dita_d5ae397b-5aa5-4de0-82c1-a4aff63c5ba1).",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "The API token for the user. This API token has no expiry; to re-generate it, delete the resource and recreate it.",
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	tflog.Trace(ctx, "create API token resource for user")

	// 1. read terraform plan data into model
	var planData ApiTokenResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 2. create resource & fill model data
	generateApiTokenInp := user.NewGenerateApiTokenInput(planData.Username.ValueString())
	generateApiTokenOutp, err := r.client.GenerateApiToken(ctx, *generateApiTokenInp)
	if err != nil {
		res.Diagnostics.AddError("failed to geneerate API token for user", err.Error())
		return
	}
	planData.ApiToken = types.StringValue(generateApiTokenOutp.ApiToken)

	// 3. fill terraform state using model data
	res.Diagnostics.Append(res.State.Set(ctx, &planData)...)
	tflog.Trace(ctx, "create user resource done")
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	tflog.Trace(ctx, "Nothing to do to update an API token")
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	tflog.Trace(ctx, "Revoke user API token")

	// 1. read state data from terraform state
	var stateData ApiTokenResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// 2. delete the resource
	deleteUserInput := user.RevokeApiTokenInput{
		Name: stateData.Username.ValueString(),
	}
	_, err := r.client.RevokeApiToken(ctx, deleteUserInput)
	if err != nil {
		res.Diagnostics.AddError("failed to delete User resource", err.Error())
	}
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "cannot read an API token once generated")
}
