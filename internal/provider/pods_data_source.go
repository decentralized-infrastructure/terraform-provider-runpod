package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &PodsDataSource{}

func NewPodsDataSource() datasource.DataSource {
	return &PodsDataSource{}
}

type PodsDataSource struct {
	client *Client
}

type PodsDataSourceModel struct {
	Pods []PodDataModel `tfsdk:"pods"`
}

type PodDataModel struct {
	ID                types.String  `tfsdk:"id"`
	Name              types.String  `tfsdk:"name"`
	ImageName         types.String  `tfsdk:"image_name"`
	DesiredStatus     types.String  `tfsdk:"desired_status"`
	PublicIp          types.String  `tfsdk:"public_ip"`
	MachineId         types.String  `tfsdk:"machine_id"`
	CostPerHr         types.Float64 `tfsdk:"cost_per_hr"`
	AdjustedCostPerHr types.Float64 `tfsdk:"adjusted_cost_per_hr"`
	MemoryInGb        types.Float64 `tfsdk:"memory_in_gb"`
	VCPUCount         types.Float64 `tfsdk:"vcpu_count"`
	GPUCount          types.Int64   `tfsdk:"gpu_count"`
}

func (d *PodsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pods"
}

func (d *PodsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source to list all RunPod Pods.",

		Attributes: map[string]schema.Attribute{
			"pods": schema.ListNestedAttribute{
				MarkdownDescription: "List of Pods.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier of the Pod.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the Pod.",
							Computed:            true,
						},
						"image_name": schema.StringAttribute{
							MarkdownDescription: "The Docker image name.",
							Computed:            true,
						},
						"desired_status": schema.StringAttribute{
							MarkdownDescription: "The desired status of the Pod.",
							Computed:            true,
						},
						"public_ip": schema.StringAttribute{
							MarkdownDescription: "The public IP address of the Pod.",
							Computed:            true,
						},
						"machine_id": schema.StringAttribute{
							MarkdownDescription: "The machine ID where the Pod is running.",
							Computed:            true,
						},
						"cost_per_hr": schema.Float64Attribute{
							MarkdownDescription: "The cost per hour in RunPod credits.",
							Computed:            true,
						},
						"adjusted_cost_per_hr": schema.Float64Attribute{
							MarkdownDescription: "The adjusted cost per hour with savings plans applied.",
							Computed:            true,
						},
						"memory_in_gb": schema.Float64Attribute{
							MarkdownDescription: "The amount of RAM in GB.",
							Computed:            true,
						},
						"vcpu_count": schema.Float64Attribute{
							MarkdownDescription: "The number of vCPUs.",
							Computed:            true,
						},
						"gpu_count": schema.Int64Attribute{
							MarkdownDescription: "The number of GPUs.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *PodsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *PodsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PodsDataSourceModel

	// Initialize empty slice
	data.Pods = []PodDataModel{}

	tflog.Debug(ctx, "Reading Pods data source")

	pods, err := d.client.ListPods(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list pods, got error: %s", err))
		return
	}

	for _, pod := range pods {
		podData := PodDataModel{
			ID:                types.StringValue(pod.ID),
			Name:              types.StringValue(pod.Name),
			ImageName:         types.StringValue(pod.ImageName),
			DesiredStatus:     types.StringValue(pod.DesiredStatus),
			PublicIp:          types.StringValue(pod.PublicIp),
			MachineId:         types.StringValue(pod.MachineId),
			CostPerHr:         types.Float64Value(pod.CostPerHr),
			AdjustedCostPerHr: types.Float64Value(pod.AdjustedCostPerHr),
			MemoryInGb:        types.Float64Value(pod.MemoryInGb),
			VCPUCount:         types.Float64Value(float64(pod.VCPUCount)),
			GPUCount:          types.Int64Value(int64(pod.GPUCount)),
		}
		data.Pods = append(data.Pods, podData)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
