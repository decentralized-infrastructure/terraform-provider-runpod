package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &NetworkVolumeResource{}
var _ resource.ResourceWithImportState = &NetworkVolumeResource{}

func NewNetworkVolumeResource() resource.Resource {
	return &NetworkVolumeResource{}
}

type NetworkVolumeResource struct {
	client *Client
}

type NetworkVolumeResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Size         types.Int64  `tfsdk:"size"`
	DataCenterId types.String `tfsdk:"data_center_id"`
}

func (r *NetworkVolumeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_volume"
}

func (r *NetworkVolumeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "RunPod Network Volume resource. A Network Volume is persistent storage that can be attached to Pods and Endpoints.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the Network Volume.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A user-defined name for the Network Volume. The name does not need to be unique.",
				Required:            true,
			},
			"size": schema.Int64Attribute{
				MarkdownDescription: "The amount of disk space, in gigabytes (GB), allocated to the Network Volume. Must be between 0 and 4000.",
				Required:            true,
			},
			"data_center_id": schema.StringAttribute{
				MarkdownDescription: "The RunPod data center ID where the Network Volume is located.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *NetworkVolumeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *NetworkVolumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NetworkVolumeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Network Volume")

	input := &NetworkVolumeCreateInput{
		Name:         data.Name.ValueString(),
		Size:         int(data.Size.ValueInt64()),
		DataCenterId: data.DataCenterId.ValueString(),
	}

	volume, err := r.client.CreateNetworkVolume(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create network volume, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Created Network Volume", map[string]interface{}{"id": volume.ID})

	data.ID = types.StringValue(volume.ID)
	data.Name = types.StringValue(volume.Name)
	data.Size = types.Int64Value(int64(volume.Size))
	data.DataCenterId = types.StringValue(volume.DataCenterId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetworkVolumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetworkVolumeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Network Volume", map[string]interface{}{"id": data.ID.ValueString()})

	volume, err := r.client.GetNetworkVolume(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network volume, got error: %s", err))
		return
	}

	data.ID = types.StringValue(volume.ID)
	data.Name = types.StringValue(volume.Name)
	data.Size = types.Int64Value(int64(volume.Size))
	data.DataCenterId = types.StringValue(volume.DataCenterId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetworkVolumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NetworkVolumeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating Network Volume", map[string]interface{}{"id": data.ID.ValueString()})

	input := &NetworkVolumeUpdateInput{
		Name: data.Name.ValueString(),
	}

	if !data.Size.IsNull() {
		size := int(data.Size.ValueInt64())
		input.Size = &size
	}

	volume, err := r.client.UpdateNetworkVolume(ctx, data.ID.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update network volume, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Updated Network Volume", map[string]interface{}{"id": volume.ID})

	data.ID = types.StringValue(volume.ID)
	data.Name = types.StringValue(volume.Name)
	data.Size = types.Int64Value(int64(volume.Size))
	data.DataCenterId = types.StringValue(volume.DataCenterId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetworkVolumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NetworkVolumeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting Network Volume", map[string]interface{}{"id": data.ID.ValueString()})

	err := r.client.DeleteNetworkVolume(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete network volume, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Deleted Network Volume", map[string]interface{}{"id": data.ID.ValueString()})
}

func (r *NetworkVolumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
