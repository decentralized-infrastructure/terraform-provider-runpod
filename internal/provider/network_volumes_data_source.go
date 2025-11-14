package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &NetworkVolumesDataSource{}

func NewNetworkVolumesDataSource() datasource.DataSource {
	return &NetworkVolumesDataSource{}
}

type NetworkVolumesDataSource struct {
	client *Client
}

type NetworkVolumesDataSourceModel struct {
	NetworkVolumes []NetworkVolumeDataModel `tfsdk:"network_volumes"`
}

type NetworkVolumeDataModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Size         types.Int64  `tfsdk:"size"`
	DataCenterId types.String `tfsdk:"data_center_id"`
}

func (d *NetworkVolumesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_volumes"
}

func (d *NetworkVolumesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source to list all RunPod Network Volumes.",

		Attributes: map[string]schema.Attribute{
			"network_volumes": schema.ListNestedAttribute{
				MarkdownDescription: "List of Network Volumes.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier of the Network Volume.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the Network Volume.",
							Computed:            true,
						},
						"size": schema.Int64Attribute{
							MarkdownDescription: "The size of the Network Volume in GB.",
							Computed:            true,
						},
						"data_center_id": schema.StringAttribute{
							MarkdownDescription: "The data center ID where the Network Volume is located.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *NetworkVolumesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworkVolumesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NetworkVolumesDataSourceModel

	// Initialize empty slice
	data.NetworkVolumes = []NetworkVolumeDataModel{}

	tflog.Debug(ctx, "Reading Network Volumes data source")

	volumes, err := d.client.ListNetworkVolumes(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list network volumes, got error: %s", err))
		return
	}

	for _, vol := range volumes {
		volumeData := NetworkVolumeDataModel{
			ID:           types.StringValue(vol.ID),
			Name:         types.StringValue(vol.Name),
			Size:         types.Int64Value(int64(vol.Size)),
			DataCenterId: types.StringValue(vol.DataCenterId),
		}
		data.NetworkVolumes = append(data.NetworkVolumes, volumeData)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
