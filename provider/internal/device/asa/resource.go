package asa

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/publicapilabels"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"

	"github.com/CiscoDevnet/terraform-provider-cdo/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &AsaDeviceResource{}
var _ resource.ResourceWithImportState = &AsaDeviceResource{}

func NewAsaDeviceResource() resource.Resource {
	return &AsaDeviceResource{}
}

type AsaDeviceResource struct {
	client *cdoClient.Client
}

type AsaDeviceResourceModel struct {
	ID            types.String `tfsdk:"id"`
	ConnectorType types.String `tfsdk:"connector_type"`
	ConnectorName types.String `tfsdk:"connector_name"`
	Name          types.String `tfsdk:"name"`
	SocketAddress types.String `tfsdk:"socket_address"`
	Host          types.String `tfsdk:"host"`
	Port          types.Int64  `tfsdk:"port"`
	Labels        types.Set    `tfsdk:"labels"`
	GroupedLabels types.Map    `tfsdk:"grouped_labels"`

	Username          types.String `tfsdk:"username"`
	Password          types.String `tfsdk:"password"`
	IgnoreCertificate types.Bool   `tfsdk:"ignore_certificate"`
	SoftwareVersion   types.String `tfsdk:"software_version"`
	AsdmVersion       types.String `tfsdk:"asdm_version"`
}

func (r *AsaDeviceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asa_device"
}

func (r *AsaDeviceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Provides an ASA device resource. This allows ASA devices to be onboarded, updated, and deleted on CDO.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the device. This is a UUID and is automatically generated when the device is created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A human-readable name for the device.",
				Required:            true,
			},
			"connector_name": schema.StringAttribute{
				MarkdownDescription: "The name of the Secure Device Connector (SDC) that will be used to communicate with the device. This value is not required if the connector type selected is Cloud Connector (CDG).",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"connector_type": schema.StringAttribute{
				MarkdownDescription: "The type of the connector that will be used to communicate with the device. CDO can communicate with your device using either a Cloud Connector (CDG) or a Secure Device Connector (SDC); see [the CDO documentation](https://docs.defenseorchestrator.com/c-connect-cisco-defense-orchestratortor-the-secure-device-connector.html) to learn more (Valid values: [CDG, SDC]).",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("CDG", "SDC"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"socket_address": schema.StringAttribute{
				MarkdownDescription: "The address of the device to onboard, specified in the format `host:port`.",
				Required:            true,
				Validators: []validator.String{
					validators.ValidateSocketAddress(),
				},
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The port used to connect to the device.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The host used to connect to the device.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"labels": schema.SetAttribute{
				MarkdownDescription: "Specify a set of labels to identify the device as part of a group. Refer to the [CDO documentation](https://docs.defenseorchestrator.com/t-applying-labels-to-devices-and-objects.html#!c-labels-and-filtering.html) for details on how labels are used in CDO.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})), // default to empty list
			},
			"grouped_labels": schema.MapAttribute{
				MarkdownDescription: "Specify a map of grouped labels to identify the device as part of a group. Refer to the [CDO documentation](https://docs.defenseorchestrator.com/t-applying-labels-to-devices-and-objects.html#!c-labels-and-filtering.html) for details on how labels are used in CDO.",
				Optional:            true,
				Computed:            true,
				ElementType: types.SetType{
					ElemType: types.StringType,
				},
				Default: mapdefault.StaticValue(types.MapValueMust(types.SetType{ElemType: types.StringType}, map[string]attr.Value{})), // default to empty list
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username used to authenticate with the device.",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password used to authenticate with the device.",
				Required:            true,
				Sensitive:           true,
			},
			"ignore_certificate": schema.BoolAttribute{
				MarkdownDescription: "Set this attribute to true if you do not want CDO to validate the certificate of this device before onboarding.",
				Required:            true,
			},
			"software_version": schema.StringAttribute{
				MarkdownDescription: "The version of the ASA device. If this attribute is set during resource creation and the version of the ASA is not the same as that specified, resource creation will fail. If the version attribute is updated following the creation of a resource, the CDO terraform provider will attempt to upgrade the device to the specified version.",
				Optional:            true,
				Computed:            true,
			},
			"asdm_version": schema.StringAttribute{
				MarkdownDescription: "The version of the ASDM on the ASA device. If this attribute is set during resource creation and the version of ASDM on the ASA is not the same as that specified, resource creation will fail. If the version attribute is updated following the creation of a resource, the CDO terraform provider will attempt to upgrade the ASDM on the device to the specified version.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *AsaDeviceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

// TODO plan diffing should exclude host, id, port, and sdc_name (unless SDCType is changed, in which case it should be a destroy and create)
// TODO terraform should wait when credentials are updated and the device is synced
// TODO verify changing groups of changes

func (r *AsaDeviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	tflog.Trace(ctx, "read ASA device resource")

	var stateData *AsaDeviceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// read asa
	readInp := asa.ReadInput{
		Uid: stateData.ID.ValueString(),
	}
	asaReadOutp, err := r.client.ReadAsa(ctx, readInp)
	if err != nil {
		if util.Is404Error(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to read ASA Device", err.Error())
		return
	}

	port, err := strconv.ParseInt(asaReadOutp.Port, 10, 16)
	if err != nil {
		resp.Diagnostics.AddError("unable to read ASA Device", err.Error())
		return
	}
	stateData.Port = types.Int64Value(port)

	stateData.ID = types.StringValue(asaReadOutp.Uid)
	stateData.ConnectorType = types.StringValue(asaReadOutp.ConnectorType)
	stateData.Name = types.StringValue(asaReadOutp.Name)
	stateData.SocketAddress = types.StringValue(asaReadOutp.SocketAddress)
	stateData.Host = types.StringValue(asaReadOutp.Host)
	stateData.IgnoreCertificate = types.BoolValue(asaReadOutp.IgnoreCertificate)
	stateData.Labels = util.GoStringSliceToTFStringSet(asaReadOutp.Tags.UngroupedTags())
	stateData.GroupedLabels = util.GoMapToStringSetTFMap(asaReadOutp.Tags.GroupedTags())
	stateData.SoftwareVersion = types.StringValue(asaReadOutp.SoftwareVersion)
	stateData.AsdmVersion = types.StringValue(asaReadOutp.AsdmVersion)

	tflog.Trace(ctx, "done read ASA device resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}

func (r *AsaDeviceResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {

	tflog.Trace(ctx, "create ASA device resource")

	var planData AsaDeviceResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if res.Diagnostics.HasError() {
		return
	}

	var specificSdcOutp *connector.ReadOutput
	if strings.EqualFold(planData.ConnectorType.ValueString(), "SDC") {
		readSdcByNameInp := connector.NewReadByNameInput(
			planData.ConnectorName.ValueString(),
		)

		var err error
		specificSdcOutp, err = r.client.ReadConnectorByName(ctx, *readSdcByNameInp)
		if err != nil {
			res.Diagnostics.AddError("failed to create ASA", err.Error())
			return
		}

	} else {
		specificSdcOutp = &connector.ReadOutput{}
	}

	// convert tf tags to go tags
	planTags, err := labelsFromAsaDeviceResourceModel(ctx, &planData)
	if err != nil {
		res.Diagnostics.AddError("error while converting terraform tags to go tags", err.Error())
		return
	}

	createInp := asa.NewCreateRequestInput(
		planData.Name.ValueString(),
		specificSdcOutp.Uid,
		planData.ConnectorType.ValueString(),
		planData.SocketAddress.ValueString(),
		planData.Username.ValueString(),
		planData.Password.ValueString(),
		planData.IgnoreCertificate.ValueBool(),
		planTags,
		planData.SoftwareVersion.ValueString(),
		planData.AsdmVersion.ValueString(),
	)

	createOutp, createSpecificOutp, createErr := r.client.CreateAsa(ctx, *createInp)
	if createErr != nil {
		tflog.Error(ctx, "Failed to create ASA device")
		if createErr.CreatedResourceId != nil {
			deleteInp := asa.NewDeleteInput(*createErr.CreatedResourceId)
			_, err := r.client.DeleteAsa(ctx, *deleteInp)
			if err != nil {
				res.Diagnostics.AddError("failed to delete ASA device", err.Error())
			}
		}

		res.Diagnostics.AddError("failed to create ASA device", createErr.Error())
		return
	}

	planData.ID = types.StringValue(createOutp.Uid)
	planData.ConnectorType = types.StringValue(createOutp.ConnectorType)
	planData.ConnectorName = getConnectorName(&planData)
	planData.Name = types.StringValue(createOutp.Name)
	parts := strings.Split(planData.SocketAddress.ValueString(), ":")
	if len(parts) == 2 {
		planData.Host = types.StringValue(parts[0])
		port, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			res.Diagnostics.AddError("failed to parse port", err.Error())
			return
		}
		planData.Port = types.Int64Value(port)
	} else {
		res.Diagnostics.AddError("invalid socket address format", "expected format is host:port")
		return
	}

	planData.SoftwareVersion = types.StringValue(createOutp.SoftwareVersion)
	planData.AsdmVersion = types.StringValue(createSpecificOutp.Metadata.AsdmVersion)

	res.Diagnostics.Append(res.State.Set(ctx, &planData)...)
}

func (r *AsaDeviceResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {

	tflog.Trace(ctx, "update ASA device resource")

	var planData *AsaDeviceResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if res.Diagnostics.HasError() {
		return
	}

	var stateData *AsaDeviceResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if res.Diagnostics.HasError() {
		return
	}

	// convert tf tags to go tags
	planTags, err := tagsFromAsaDeviceResourceModel(ctx, planData)
	if err != nil {
		res.Diagnostics.AddError("error while converting terraform tags to go tags", err.Error())
		return
	}

	updateInp := asa.NewUpdateInput(
		stateData.ID.ValueString(),
		stateData.Name.ValueString(),
		"",
		"",
		planTags,
	)

	if isNameUpdated(planData, stateData) {
		updateInp.Name = planData.Name.ValueString()
	}

	if isLocationUpdated(planData, stateData) {
		updateInp.Location = planData.SocketAddress.ValueString()
	}

	if isCredentialUpdated(planData, stateData) {
		stateData.Username = types.StringValue(planData.Username.ValueString())
		stateData.Password = types.StringValue(planData.Password.ValueString())
		res.Diagnostics.Append(res.State.Set(ctx, &stateData)...)
		if res.Diagnostics.HasError() {
			return
		}

		updateInp.Username = planData.Username.ValueString()
		updateInp.Password = planData.Password.ValueString()
	}

	updateOutp, err := r.client.UpdateAsa(ctx, *updateInp)
	if err != nil {
		res.Diagnostics.AddError("failed to update ASA device", err.Error())
		return
	}

	port, err := parsePort(updateOutp.Port)
	if err != nil {
		res.Diagnostics.AddError("unable to parse port", err.Error())
		return
	}

	stateData.ID = types.StringValue(updateOutp.Uid)
	stateData.ConnectorType = types.StringValue(planData.ConnectorType.ValueString())
	stateData.ConnectorName = getConnectorName(planData)
	stateData.Name = types.StringValue(updateOutp.Name)
	stateData.SocketAddress = planData.SocketAddress
	stateData.Host = types.StringValue(updateOutp.Host)
	stateData.Port = types.Int64Value(port)
	stateData.Labels = planData.Labels
	stateData.GroupedLabels = planData.GroupedLabels

	stateData.IgnoreCertificate = planData.IgnoreCertificate

	res.Diagnostics.Append(res.State.Set(ctx, &stateData)...)
}

func (r *AsaDeviceResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {

	tflog.Trace(ctx, "delete ASA device resource")

	var stateData AsaDeviceResourceModel

	res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if res.Diagnostics.HasError() {
		return
	}

	deleteInp := asa.NewDeleteInput(stateData.ID.ValueString())
	_, err := r.client.DeleteAsa(ctx, *deleteInp)
	if err != nil {
		res.Diagnostics.AddError("failed to delete ASA device", err.Error())
		return
	}

}

func (r *AsaDeviceResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, res *resource.ModifyPlanResponse) {
	if !req.State.Raw.IsNull() {
		// this is an update
		var stateData *AsaDeviceResourceModel
		res.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
		if res.Diagnostics.HasError() {
			return
		}

		var planData *AsaDeviceResourceModel
		res.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
		if res.Diagnostics.HasError() {
			return
		}

		if planData != nil && stateData != nil && strings.EqualFold(planData.SocketAddress.ValueString(), stateData.SocketAddress.ValueString()) {
			tflog.Debug(ctx, "There is no change in the socket address; remove host and port diffs")
			planData.Host = stateData.Host
			planData.Port = stateData.Port
		}

		res.Diagnostics.Append(res.Plan.Set(ctx, &planData)...)
	}
}

func (r *AsaDeviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}

func isCredentialUpdated(planData, stateData *AsaDeviceResourceModel) bool {
	return planData.Username.ValueString() != stateData.Username.ValueString() || planData.Password.ValueString() != stateData.Password.ValueString()
}

func isNameUpdated(planData, stateData *AsaDeviceResourceModel) bool {
	return !planData.Name.Equal(stateData.Name)
}

func isLocationUpdated(planData, stateData *AsaDeviceResourceModel) bool {
	return !planData.SocketAddress.Equal(stateData.SocketAddress)
}

func parsePort(rawPort string) (int64, error) {
	return strconv.ParseInt(rawPort, 10, 16)

}

func getConnectorName(planData *AsaDeviceResourceModel) basetypes.StringValue {
	if planData.ConnectorName.ValueString() != "" {
		return types.StringValue(planData.ConnectorName.ValueString())
	} else {
		return types.StringNull()
	}
}

func ungroupedAndGroupedLabelsFromResourceModel(ctx context.Context, resourceModel *AsaDeviceResourceModel) ([]string, map[string][]string, error) {
	if resourceModel == nil {
		return nil, nil, errors.New("resource model cannot be nil")
	}

	ungroupedLabels, err := util.TFStringSetToGoStringList(ctx, resourceModel.Labels)
	if err != nil {
		return nil, nil, fmt.Errorf("error while converting terraform labels to go slice, %s", resourceModel.Labels)
	}

	groupedLabels, err := util.TFMapToGoMapOfStringSlices(ctx, resourceModel.GroupedLabels)
	if err != nil {
		return nil, nil, fmt.Errorf("error while converting terraform grouped labels to go map, %v", resourceModel.GroupedLabels)
	}

	return ungroupedLabels, groupedLabels, nil
}

func tagsFromAsaDeviceResourceModel(ctx context.Context, resourceModel *AsaDeviceResourceModel) (tags.Type, error) {
	ungroupedLabels, groupedLabels, err := ungroupedAndGroupedLabelsFromResourceModel(ctx, resourceModel)
	if err != nil {
		return nil, err
	}

	return tags.New(ungroupedLabels, groupedLabels), nil
}

func labelsFromAsaDeviceResourceModel(ctx context.Context, resourceModel *AsaDeviceResourceModel) (publicapilabels.Type, error) {
	ungroupedLabels, groupedLabels, err := ungroupedAndGroupedLabelsFromResourceModel(ctx, resourceModel)
	if err != nil {
		return publicapilabels.Empty(), err
	}

	return publicapilabels.New(ungroupedLabels, groupedLabels), nil
}
