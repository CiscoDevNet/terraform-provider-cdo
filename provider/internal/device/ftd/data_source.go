package ftd

import (
	"context"
	"fmt"
	cdoClient "github.com/CiscoDevnet/terraform-provider-cdo/go-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = NewDataSource()

// NewDataSource is used in provider.go to register this data source with terraform.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource is the struct object that will be consumed by terraform, it contains methods that defines metadata, schema, read, create, etc...
type DataSource struct {
	client *cdoClient.Client
}

/////
// model class: define mapping from Go types to Terraform types.
/////

type DataSourceModel struct {
	ID               types.String   `tfsdk:"id"`
	Name             types.String   `tfsdk:"name"`
	AccessPolicyName types.String   `tfsdk:"access_policy_name"`
	PerformanceTier  types.String   `tfsdk:"performance_tier"`
	Virtual          types.Bool     `tfsdk:"virtual"`
	Licenses         []types.String `tfsdk:"licenses"`
	Labels           []types.String `tfsdk:"labels"`

	AccessPolicyUid  types.String `tfsdk:"access_policy_id"`
	GeneratedCommand types.String `tfsdk:"generated_command"`
	Hostname         types.String `tfsdk:"hostname"`
	NatId            types.String `tfsdk:"nat_id"`
	RegKey           types.String `tfsdk:"reg_key"`
}

// Metadata is primarily used to define the name for this data source.
func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ios_device"
}

// Schema is primarily used to define the terraform schema for this data source.
// it needs to be consistent with the Model classes' `tfsdk:"xxx"` above.
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Ftd data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the device. This is a UUID and is automatically generated when the device is created.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A human-readable name for the Firewall Threat Defense (FTD). This name must be unique.",
				Required:            true,
			},
			"access_policy_name": schema.StringAttribute{
				MarkdownDescription: "The name of the Cloud-Delivered FMC (cdFMC) access policy that will be used by the FTD.",
				Computed:            true,
			},
			"performance_tier": schema.StringAttribute{
				MarkdownDescription: "The performance tier of the virtual FTD, if virtual is set to false, this field is ignored as performance tiers are not applicable to physical FTD devices. Allowed values are: [\"FTDv5\", \"FTDv10\", \"FTDv20\", \"FTDv30\", \"FTDv50\", \"FTDv100\", \"FTDv\"].",
				Computed:            true,
			},
			"virtual": schema.BoolAttribute{
				MarkdownDescription: "This determines if this FTD is virtual. If false, performance_tier is ignored as performance tiers are not applicable to physical FTD devices.",
				Computed:            true,
			},
			"licenses": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Comma-separated list of licenses to apply to this FTD. You must enable at least the \"BASE\" license. Allowed values are: [\"BASE\", \"CARRIER\", \"THREAT\", \"MALWARE\", \"URLFilter\",].",
				Computed:            true,
			},
			"labels": schema.ListAttribute{
				MarkdownDescription: "Set a list of labels to identify the device as part of a group. Refer to the [CDO documentation](https://docs.defenseorchestrator.com/t-applying-labels-to-devices-and-objects.html#!c-labels-and-filtering.html) for details on how labels are used in CDO.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"generated_command": schema.StringAttribute{
				MarkdownDescription: "The command to run in the FTD CLI to register it with the cloud-delivered FMC (cdFMC).",
				Computed:            true,
			},
			"access_policy_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the cloud-delivered FMC (cdFMC) access policy applied to this FTD.",
				Computed:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The Hostname of the cloud-delivered FMC (cdFMC) manages this FTD.",
				Computed:            true,
			},
			"nat_id": schema.StringAttribute{
				MarkdownDescription: "The Network Address Translation (NAT) ID of this FTD.",
				Computed:            true,
			},
			"reg_key": schema.StringAttribute{
				MarkdownDescription: "The Registration Key of this FTD.",
				Computed:            true,
			},
		},
	}
}

// Configure sets the environment for CRUD operations, i.e. we set the CDO Go client to use when doing the CRUD operations.
// The Go client is configured in the provider, and it is available here, so we just set it directly.
func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cdoClient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *cdoClient.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Read operation is all it needs for terraform data source.
// this function is responsible for:
// 1. mapping the cdo go client's data to the model classes we defined above.
// 2. report any error using `resp.Diagnostics`.
// 3. set model classes to the terraform state.
func (r *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	tflog.Trace(ctx, "reading Ftd data source")

	// 1. read state data
	var stateData DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// 2. do read
	if err := ReadDataSource(ctx, r, &stateData); err != nil {
		resp.Diagnostics.AddError("failed to read Ftd data source", err.Error())
	}

	// 3. save data into terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
}
