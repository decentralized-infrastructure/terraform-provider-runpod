package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &EndpointResource{}
var _ resource.ResourceWithImportState = &EndpointResource{}

func NewEndpointResource() resource.Resource {
	return &EndpointResource{}
}

type EndpointResource struct {
	client *Client
}

type EndpointResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	TemplateId          types.String `tfsdk:"template_id"`
	ComputeType         types.String `tfsdk:"compute_type"`
	GPUCount            types.Int64  `tfsdk:"gpu_count"`
	VCPUCount           types.Int64  `tfsdk:"vcpu_count"`
	GPUTypeIds          types.List   `tfsdk:"gpu_type_ids"`
	CPUFlavorIds        types.List   `tfsdk:"cpu_flavor_ids"`
	DataCenterIds       types.List   `tfsdk:"data_center_ids"`
	NetworkVolumeId     types.String `tfsdk:"network_volume_id"`
	WorkersMin          types.Int64  `tfsdk:"workers_min"`
	WorkersMax          types.Int64  `tfsdk:"workers_max"`
	IdleTimeout         types.Int64  `tfsdk:"idle_timeout"`
	ExecutionTimeoutMs  types.Int64  `tfsdk:"execution_timeout_ms"`
	ScalerType          types.String `tfsdk:"scaler_type"`
	ScalerValue         types.Int64  `tfsdk:"scaler_value"`
	AllowedCudaVersions types.List   `tfsdk:"allowed_cuda_versions"`
	Flashboot           types.Bool   `tfsdk:"flashboot"`
	// Computed fields
	CreatedAt types.String `tfsdk:"created_at"`
	UserId    types.String `tfsdk:"user_id"`
	Version   types.Int64  `tfsdk:"version"`
}

func (r *EndpointResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_endpoint"
}

func (r *EndpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "RunPod Serverless Endpoint resource. An Endpoint is a serverless deployment that auto-scales based on demand.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the Endpoint.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A user-defined name for the Endpoint.",
				Optional:            true,
			},
			"template_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the template used to create the Endpoint.",
				Required:            true,
			},
			"compute_type": schema.StringAttribute{
				MarkdownDescription: "Set to GPU for GPU workers or CPU for CPU workers.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("GPU"),
			},
			"gpu_count": schema.Int64Attribute{
				MarkdownDescription: "The number of GPUs attached to each worker on the Endpoint.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
			"vcpu_count": schema.Int64Attribute{
				MarkdownDescription: "If the Endpoint is a CPU endpoint, the number of vCPUs allocated to each worker.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(2),
			},
			"gpu_type_ids": schema.ListAttribute{
				MarkdownDescription: "A list of RunPod GPU types which can be attached to workers.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"cpu_flavor_ids": schema.ListAttribute{
				MarkdownDescription: "A list of RunPod CPU flavors which can be attached to workers.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"data_center_ids": schema.ListAttribute{
				MarkdownDescription: "A list of RunPod data center IDs where workers can be located.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"network_volume_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the network volume to attach to the Endpoint.",
				Optional:            true,
			},
			"workers_min": schema.Int64Attribute{
				MarkdownDescription: "The minimum number of workers that will run at the same time.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
			},
			"workers_max": schema.Int64Attribute{
				MarkdownDescription: "The maximum number of workers that can be running at the same time.",
				Optional:            true,
			},
			"idle_timeout": schema.Int64Attribute{
				MarkdownDescription: "The number of seconds a worker can run without taking a job before the worker is scaled down.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(5),
			},
			"execution_timeout_ms": schema.Int64Attribute{
				MarkdownDescription: "The maximum number of milliseconds a request can run before the worker is stopped.",
				Optional:            true,
			},
			"scaler_type": schema.StringAttribute{
				MarkdownDescription: "The method used to scale up workers. QUEUE_DELAY or REQUEST_COUNT.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("QUEUE_DELAY"),
			},
			"scaler_value": schema.Int64Attribute{
				MarkdownDescription: "For QUEUE_DELAY: seconds a request can remain in queue. For REQUEST_COUNT: target requests per worker.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(4),
			},
			"allowed_cuda_versions": schema.ListAttribute{
				MarkdownDescription: "A list of acceptable CUDA versions on the workers.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"flashboot": schema.BoolAttribute{
				MarkdownDescription: "Whether to use flash boot for the Endpoint.",
				Optional:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The UTC timestamp when the Endpoint was created.",
				Computed:            true,
			},
			"user_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the user who created the Endpoint.",
				Computed:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "The version number of the Endpoint.",
				Computed:            true,
			},
		},
	}
}

func (r *EndpointResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EndpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Endpoint")

	input := &EndpointCreateInput{
		Name:            data.Name.ValueString(),
		TemplateId:      data.TemplateId.ValueString(),
		ComputeType:     data.ComputeType.ValueString(),
		NetworkVolumeId: data.NetworkVolumeId.ValueString(),
		ScalerType:      data.ScalerType.ValueString(),
	}

	if !data.GPUCount.IsNull() {
		gpuCount := int(data.GPUCount.ValueInt64())
		input.GPUCount = &gpuCount
	}
	if !data.VCPUCount.IsNull() {
		vcpuCount := int(data.VCPUCount.ValueInt64())
		input.VCPUCount = &vcpuCount
	}
	if !data.WorkersMin.IsNull() {
		workersMin := int(data.WorkersMin.ValueInt64())
		input.WorkersMin = &workersMin
	}
	if !data.WorkersMax.IsNull() {
		workersMax := int(data.WorkersMax.ValueInt64())
		input.WorkersMax = &workersMax
	}
	if !data.IdleTimeout.IsNull() {
		idleTimeout := int(data.IdleTimeout.ValueInt64())
		input.IdleTimeout = &idleTimeout
	}
	if !data.ExecutionTimeoutMs.IsNull() {
		execTimeout := int(data.ExecutionTimeoutMs.ValueInt64())
		input.ExecutionTimeoutMs = &execTimeout
	}
	if !data.ScalerValue.IsNull() {
		scalerValue := int(data.ScalerValue.ValueInt64())
		input.ScalerValue = &scalerValue
	}
	if !data.Flashboot.IsNull() {
		flashboot := data.Flashboot.ValueBool()
		input.Flashboot = &flashboot
	}

	if !data.GPUTypeIds.IsNull() {
		resp.Diagnostics.Append(data.GPUTypeIds.ElementsAs(ctx, &input.GPUTypeIds, false)...)
	}
	if !data.CPUFlavorIds.IsNull() {
		resp.Diagnostics.Append(data.CPUFlavorIds.ElementsAs(ctx, &input.CPUFlavorIds, false)...)
	}
	if !data.DataCenterIds.IsNull() {
		resp.Diagnostics.Append(data.DataCenterIds.ElementsAs(ctx, &input.DataCenterIds, false)...)
	}
	if !data.AllowedCudaVersions.IsNull() {
		resp.Diagnostics.Append(data.AllowedCudaVersions.ElementsAs(ctx, &input.AllowedCudaVersions, false)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint, err := r.client.CreateEndpoint(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create endpoint, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Created Endpoint", map[string]interface{}{"id": endpoint.ID})

	r.updateStateFromEndpoint(ctx, &data, endpoint)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EndpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Endpoint", map[string]interface{}{"id": data.ID.ValueString()})

	endpoint, err := r.client.GetEndpoint(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read endpoint, got error: %s", err))
		return
	}

	r.updateStateFromEndpoint(ctx, &data, endpoint)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EndpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating Endpoint", map[string]interface{}{"id": data.ID.ValueString()})

	input := &EndpointUpdateInput{
		Name:            data.Name.ValueString(),
		TemplateId:      data.TemplateId.ValueString(),
		NetworkVolumeId: data.NetworkVolumeId.ValueString(),
		ScalerType:      data.ScalerType.ValueString(),
	}

	if !data.GPUCount.IsNull() {
		gpuCount := int(data.GPUCount.ValueInt64())
		input.GPUCount = &gpuCount
	}
	if !data.VCPUCount.IsNull() {
		vcpuCount := int(data.VCPUCount.ValueInt64())
		input.VCPUCount = &vcpuCount
	}
	if !data.WorkersMin.IsNull() {
		workersMin := int(data.WorkersMin.ValueInt64())
		input.WorkersMin = &workersMin
	}
	if !data.WorkersMax.IsNull() {
		workersMax := int(data.WorkersMax.ValueInt64())
		input.WorkersMax = &workersMax
	}
	if !data.IdleTimeout.IsNull() {
		idleTimeout := int(data.IdleTimeout.ValueInt64())
		input.IdleTimeout = &idleTimeout
	}
	if !data.ExecutionTimeoutMs.IsNull() {
		execTimeout := int(data.ExecutionTimeoutMs.ValueInt64())
		input.ExecutionTimeoutMs = &execTimeout
	}
	if !data.ScalerValue.IsNull() {
		scalerValue := int(data.ScalerValue.ValueInt64())
		input.ScalerValue = &scalerValue
	}
	if !data.Flashboot.IsNull() {
		flashboot := data.Flashboot.ValueBool()
		input.Flashboot = &flashboot
	}

	if !data.GPUTypeIds.IsNull() {
		resp.Diagnostics.Append(data.GPUTypeIds.ElementsAs(ctx, &input.GPUTypeIds, false)...)
	}
	if !data.CPUFlavorIds.IsNull() {
		resp.Diagnostics.Append(data.CPUFlavorIds.ElementsAs(ctx, &input.CPUFlavorIds, false)...)
	}
	if !data.DataCenterIds.IsNull() {
		resp.Diagnostics.Append(data.DataCenterIds.ElementsAs(ctx, &input.DataCenterIds, false)...)
	}
	if !data.AllowedCudaVersions.IsNull() {
		resp.Diagnostics.Append(data.AllowedCudaVersions.ElementsAs(ctx, &input.AllowedCudaVersions, false)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint, err := r.client.UpdateEndpoint(ctx, data.ID.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update endpoint, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Updated Endpoint", map[string]interface{}{"id": endpoint.ID})

	r.updateStateFromEndpoint(ctx, &data, endpoint)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EndpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting Endpoint", map[string]interface{}{"id": data.ID.ValueString()})

	err := r.client.DeleteEndpoint(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete endpoint, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Deleted Endpoint", map[string]interface{}{"id": data.ID.ValueString()})
}

func (r *EndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *EndpointResource) updateStateFromEndpoint(ctx context.Context, data *EndpointResourceModel, endpoint *Endpoint) {
	data.ID = types.StringValue(endpoint.ID)
	data.Name = types.StringValue(endpoint.Name)
	data.TemplateId = types.StringValue(endpoint.TemplateId)
	data.ComputeType = types.StringValue(endpoint.ComputeType)
	data.CreatedAt = types.StringValue(endpoint.CreatedAt)
	data.UserId = types.StringValue(endpoint.UserId)
	data.Version = types.Int64Value(int64(endpoint.Version))

	if endpoint.GPUCount > 0 {
		data.GPUCount = types.Int64Value(int64(endpoint.GPUCount))
	}
	if endpoint.VCPUCount > 0 {
		data.VCPUCount = types.Int64Value(int64(endpoint.VCPUCount))
	}
	if endpoint.WorkersMin >= 0 {
		data.WorkersMin = types.Int64Value(int64(endpoint.WorkersMin))
	}
	if endpoint.WorkersMax > 0 {
		data.WorkersMax = types.Int64Value(int64(endpoint.WorkersMax))
	}
	if endpoint.IdleTimeout > 0 {
		data.IdleTimeout = types.Int64Value(int64(endpoint.IdleTimeout))
	}
	if endpoint.ExecutionTimeoutMs > 0 {
		data.ExecutionTimeoutMs = types.Int64Value(int64(endpoint.ExecutionTimeoutMs))
	}
	if endpoint.ScalerValue > 0 {
		data.ScalerValue = types.Int64Value(int64(endpoint.ScalerValue))
	}
	if endpoint.ScalerType != "" {
		data.ScalerType = types.StringValue(endpoint.ScalerType)
	}
	if endpoint.NetworkVolumeId != "" {
		data.NetworkVolumeId = types.StringValue(endpoint.NetworkVolumeId)
	}
}
