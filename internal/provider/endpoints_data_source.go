package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &EndpointsDataSource{}

func NewEndpointsDataSource() datasource.DataSource {
	return &EndpointsDataSource{}
}

type EndpointsDataSource struct {
	client *Client
}

type EndpointsDataSourceModel struct {
	Endpoints []EndpointDataModel `tfsdk:"endpoints"`
}

type EndpointDataModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	TemplateId         types.String `tfsdk:"template_id"`
	ComputeType        types.String `tfsdk:"compute_type"`
	GPUCount           types.Int64  `tfsdk:"gpu_count"`
	VCPUCount          types.Int64  `tfsdk:"vcpu_count"`
	WorkersMin         types.Int64  `tfsdk:"workers_min"`
	WorkersMax         types.Int64  `tfsdk:"workers_max"`
	IdleTimeout        types.Int64  `tfsdk:"idle_timeout"`
	ExecutionTimeoutMs types.Int64  `tfsdk:"execution_timeout_ms"`
	ScalerType         types.String `tfsdk:"scaler_type"`
	ScalerValue        types.Int64  `tfsdk:"scaler_value"`
	CreatedAt          types.String `tfsdk:"created_at"`
	Version            types.Int64  `tfsdk:"version"`
}

func (d *EndpointsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoints"
}

func (d *EndpointsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source to list all RunPod Serverless Endpoints.",

		Attributes: map[string]schema.Attribute{
			"endpoints": schema.ListNestedAttribute{
				MarkdownDescription: "List of Endpoints.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier of the Endpoint.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the Endpoint.",
							Computed:            true,
						},
						"template_id": schema.StringAttribute{
							MarkdownDescription: "The template ID used by the Endpoint.",
							Computed:            true,
						},
						"compute_type": schema.StringAttribute{
							MarkdownDescription: "The compute type (GPU or CPU).",
							Computed:            true,
						},
						"gpu_count": schema.Int64Attribute{
							MarkdownDescription: "The number of GPUs per worker.",
							Computed:            true,
						},
						"vcpu_count": schema.Int64Attribute{
							MarkdownDescription: "The number of vCPUs per worker.",
							Computed:            true,
						},
						"workers_min": schema.Int64Attribute{
							MarkdownDescription: "The minimum number of workers.",
							Computed:            true,
						},
						"workers_max": schema.Int64Attribute{
							MarkdownDescription: "The maximum number of workers.",
							Computed:            true,
						},
						"idle_timeout": schema.Int64Attribute{
							MarkdownDescription: "The idle timeout in seconds.",
							Computed:            true,
						},
						"execution_timeout_ms": schema.Int64Attribute{
							MarkdownDescription: "The execution timeout in milliseconds.",
							Computed:            true,
						},
						"scaler_type": schema.StringAttribute{
							MarkdownDescription: "The scaler type (QUEUE_DELAY or REQUEST_COUNT).",
							Computed:            true,
						},
						"scaler_value": schema.Int64Attribute{
							MarkdownDescription: "The scaler value.",
							Computed:            true,
						},
						"created_at": schema.StringAttribute{
							MarkdownDescription: "The creation timestamp.",
							Computed:            true,
						},
						"version": schema.Int64Attribute{
							MarkdownDescription: "The endpoint version.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *EndpointsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EndpointsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EndpointsDataSourceModel

	tflog.Debug(ctx, "Reading Endpoints data source")

	endpoints, err := d.client.ListEndpoints(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list endpoints, got error: %s", err))
		return
	}

	for _, endpoint := range endpoints {
		endpointData := EndpointDataModel{
			ID:                 types.StringValue(endpoint.ID),
			Name:               types.StringValue(endpoint.Name),
			TemplateId:         types.StringValue(endpoint.TemplateId),
			ComputeType:        types.StringValue(endpoint.ComputeType),
			GPUCount:           types.Int64Value(int64(endpoint.GPUCount)),
			VCPUCount:          types.Int64Value(int64(endpoint.VCPUCount)),
			WorkersMin:         types.Int64Value(int64(endpoint.WorkersMin)),
			WorkersMax:         types.Int64Value(int64(endpoint.WorkersMax)),
			IdleTimeout:        types.Int64Value(int64(endpoint.IdleTimeout)),
			ExecutionTimeoutMs: types.Int64Value(int64(endpoint.ExecutionTimeoutMs)),
			ScalerType:         types.StringValue(endpoint.ScalerType),
			ScalerValue:        types.Int64Value(int64(endpoint.ScalerValue)),
			CreatedAt:          types.StringValue(endpoint.CreatedAt),
			Version:            types.Int64Value(int64(endpoint.Version)),
		}
		data.Endpoints = append(data.Endpoints, endpointData)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
