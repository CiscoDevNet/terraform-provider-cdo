package asa

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/sdc"
	"github.com/CiscoDevnet/terraform-provider-cdo/validators"
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

	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`

	IgnoreCertificate types.Bool `tfsdk:"ignore_certificate"`
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
				MarkdownDescription: "Unique identifier of the device. This is a UUID and will be automatically generated when the device is created.",
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
				MarkdownDescription: "The type of the connector that will be used to communicate with the device. CDO can communicate with your device using either a Cloud Connector (CDG) or a Secure Device Connector (SDC); see [the CDO documentation](https://docs.defenseorchestrator.com/c-connect-cisco-defense-orchestratortor-the-secure-device-connector.html) to learn mor (Valid values: [CDG, SDC]).",
				Required:            true,
				Validators: []validator.String{
					validators.OneOf("CDG", "SDC"),
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
			},
			"host": schema.StringAttribute{
				MarkdownDescription: "The host used to connect to the device.",
				Computed:            true,
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
				MarkdownDescription: "Set this attribute to true if you do not wish for CDO to validate the certificate of this device before onboarding.",
				Required:            true,
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
	stateData.IgnoreCertificate = types.BoolValue(asaReadOutp.IgnoreCertifcate)

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

	var specificSdcOutp *sdc.ReadOutput
	if strings.EqualFold(planData.ConnectorType.ValueString(), "SDC") {
		readSdcByNameInp := sdc.NewReadByNameInput(
			planData.ConnectorName.ValueString(),
		)

		var err error
		specificSdcOutp, err = r.client.ReadSdcByName(ctx, *readSdcByNameInp)
		if err != nil {
			res.Diagnostics.AddError("failed to read SDC by name", err.Error())
			return
		}

	} else {
		specificSdcOutp = &sdc.ReadOutput{}
	}

	createInp := asa.NewCreateRequestInput(
		planData.Name.ValueString(),
		specificSdcOutp.Uid,
		planData.ConnectorType.ValueString(),
		planData.SocketAddress.ValueString(),
		planData.Username.ValueString(),
		planData.Password.ValueString(),
		planData.IgnoreCertificate.ValueBool(),
	)

	createOutp, createErr := r.client.CreateAsa(ctx, *createInp)
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
	planData.ConnectorType = types.StringValue(createOutp.SdcType)
	planData.ConnectorName = getSdcName(&planData)
	planData.Name = types.StringValue(createOutp.Name)
	planData.Host = types.StringValue(createOutp.Host)

	port, err := strconv.ParseInt(createOutp.Port, 10, 16)
	if err != nil {
		res.Diagnostics.AddError("failed to parse ASA port", err.Error())
		// delete the ASA coz we're screwed here

	}
	planData.Port = types.Int64Value(port)

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

	updateInp := asa.NewUpdateInput(stateData.ID.ValueString(), stateData.Name.ValueString(), "", "")

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
	stateData.ConnectorName = getSdcName(planData)
	stateData.Name = types.StringValue(updateOutp.Name)
	stateData.SocketAddress = planData.SocketAddress
	stateData.Host = types.StringValue(updateOutp.Host)
	stateData.Port = types.Int64Value(port)

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

func getSdcName(planData *AsaDeviceResourceModel) basetypes.StringValue {
	if planData.ConnectorName.ValueString() != "" {
		return types.StringValue(planData.ConnectorName.ValueString())
	} else {
		return types.StringNull()
	}
}
