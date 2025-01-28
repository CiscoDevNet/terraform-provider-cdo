package ftdversion

import (
	"context"
	"fmt"
	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/cloudftd"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &Resource{}

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	client *cdoClient.Client
}

type ResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	FtdUid                  types.String `tfsdk:"ftd_uid"`
	SoftwareVersion         types.String `tfsdk:"software_version"`
	SoftwareVersionOnDevice types.String `tfsdk:"software_version_on_device"`
}

func (r *Resource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ftd_device_version"
}

func (r *Resource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "Provides a resource to upgrade the software version of an FTD device." +
			" Note: The FTD device has to already have been added to the Terraform state using a " +
			"resource or a data source.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the FTD version resource",
				Computed:            true,
			},
			"ftd_uid": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the FTD device to upgrade.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"software_version": schema.StringAttribute{
				MarkdownDescription: "The software version to upgrade the FTD device to.",
				Required:            true,
			},
			"software_version_on_device": schema.StringAttribute{
				MarkdownDescription: "The software version currently on the FTD device.",
				Computed:            true,
			},
		},
	}
}

func (r *Resource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*cdoClient.Client)

	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *cdoClient.Client, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *Resource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	tflog.Debug(ctx, "Create a new FTD version resource...")

	var planData ResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &planData)...)
	if response.Diagnostics.HasError() {
		return
	}

	ftdDevice, err := r.upgrade(ctx, planData.FtdUid.ValueString(), planData.SoftwareVersion.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Failed to upgrade FTD device...", err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("FTD device upgraded successfully: %v", ftdDevice))
	planData.Id = planData.FtdUid
	planData.SoftwareVersionOnDevice = types.StringValue(ftdDevice.SoftwareVersion)
	response.Diagnostics.Append(response.State.Set(ctx, &planData)...)
}

func (r *Resource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading FTD device to update the upgrade resource...")
	var stateData ResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &stateData)...)
	if response.Diagnostics.HasError() {
		return
	}

	ftdDevice, err := cloudftd.ReadByUid(ctx, r.client.Client, cloudftd.ReadByUidInput{Uid: stateData.FtdUid.ValueString()})
	if err != nil {
		response.Diagnostics.AddError("Failed to read FTD device...", err.Error())
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("FTD device found: %v", ftdDevice))

	stateData.SoftwareVersionOnDevice = types.StringValue(ftdDevice.SoftwareVersion)
	response.Diagnostics.Append(response.State.Set(ctx, &stateData)...)
}

func (r *Resource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update a new FTD version resource...")

	var planData ResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &planData)...)
	if response.Diagnostics.HasError() {
		return
	}

	ftdDevice, err := r.upgrade(ctx, planData.FtdUid.ValueString(), planData.SoftwareVersion.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Failed to upgrade FTD device...", err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("FTD device upgraded successfully: %v", ftdDevice))
	response.Diagnostics.Append(response.State.Set(ctx, &planData)...)
}

func (r *Resource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Info(ctx, "Removing a version resource is a noop. It will not trigger a revert of the upgrade on the FTD device.")
}

func (r *Resource) upgrade(ctx context.Context, deviceUid string, softwareVersion string) (*cloudftd.FtdDevice, error) {
	ftdUpgradeService := cloudftd.NewFtdUpgradeService(ctx, &r.client.Client)
	ftdDevice, err := ftdUpgradeService.Upgrade(deviceUid, softwareVersion)
	if err != nil {
		return nil, err
	}

	return ftdDevice, nil
}
