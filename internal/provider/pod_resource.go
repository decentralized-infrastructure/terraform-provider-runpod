package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PodResource{}
var _ resource.ResourceWithImportState = &PodResource{}

func NewPodResource() resource.Resource {
	return &PodResource{}
}

// PodResource defines the resource implementation.
type PodResource struct {
	client *Client
}

// PodResourceModel describes the resource data model.
type PodResourceModel struct {
	ID                      types.String  `tfsdk:"id"`
	Name                    types.String  `tfsdk:"name"`
	ImageName               types.String  `tfsdk:"image_name"`
	ComputeType             types.String  `tfsdk:"compute_type"`
	CloudType               types.String  `tfsdk:"cloud_type"`
	GPUCount                types.Int64   `tfsdk:"gpu_count"`
	VCPUCount               types.Int64   `tfsdk:"vcpu_count"`
	GPUTypeIds              types.List    `tfsdk:"gpu_type_ids"`
	CPUFlavorIds            types.List    `tfsdk:"cpu_flavor_ids"`
	DataCenterIds           types.List    `tfsdk:"data_center_ids"`
	ContainerDiskInGb       types.Int64   `tfsdk:"container_disk_in_gb"`
	VolumeInGb              types.Int64   `tfsdk:"volume_in_gb"`
	VolumeMountPath         types.String  `tfsdk:"volume_mount_path"`
	Ports                   types.List    `tfsdk:"ports"`
	Env                     types.Map     `tfsdk:"env"`
	DockerEntrypoint        types.List    `tfsdk:"docker_entrypoint"`
	DockerStartCmd          types.List    `tfsdk:"docker_start_cmd"`
	TemplateId              types.String  `tfsdk:"template_id"`
	NetworkVolumeId         types.String  `tfsdk:"network_volume_id"`
	Interruptible           types.Bool    `tfsdk:"interruptible"`
	Locked                  types.Bool    `tfsdk:"locked"`
	MinVCPUPerGPU           types.Int64   `tfsdk:"min_vcpu_per_gpu"`
	MinRAMPerGPU            types.Int64   `tfsdk:"min_ram_per_gpu"`
	MinDownloadMbps         types.Float64 `tfsdk:"min_download_mbps"`
	MinUploadMbps           types.Float64 `tfsdk:"min_upload_mbps"`
	MinDiskBandwidthMBps    types.Float64 `tfsdk:"min_disk_bandwidth_mbps"`
	SupportPublicIp         types.Bool    `tfsdk:"support_public_ip"`
	GlobalNetworking        types.Bool    `tfsdk:"global_networking"`
	AllowedCudaVersions     types.List    `tfsdk:"allowed_cuda_versions"`
	CountryCodes            types.List    `tfsdk:"country_codes"`
	GPUTypePriority         types.String  `tfsdk:"gpu_type_priority"`
	CPUFlavorPriority       types.String  `tfsdk:"cpu_flavor_priority"`
	DataCenterPriority      types.String  `tfsdk:"data_center_priority"`
	ContainerRegistryAuthId types.String  `tfsdk:"container_registry_auth_id"`
	// Computed fields
	DesiredStatus     types.String  `tfsdk:"desired_status"`
	PublicIp          types.String  `tfsdk:"public_ip"`
	MachineId         types.String  `tfsdk:"machine_id"`
	ActualDataCenter  types.String  `tfsdk:"actual_data_center"`
	CostPerHr         types.Float64 `tfsdk:"cost_per_hr"`
	AdjustedCostPerHr types.Float64 `tfsdk:"adjusted_cost_per_hr"`
	MemoryInGb        types.Float64 `tfsdk:"memory_in_gb"`
	LastStartedAt     types.String  `tfsdk:"last_started_at"`
}

func (r *PodResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pod"
}

func (r *PodResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "RunPod Pod resource. A Pod is a container instance that can be either GPU-based or CPU-based.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the Pod.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A user-defined name for the Pod. The name does not need to be unique.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("my pod"),
			},
			"image_name": schema.StringAttribute{
				MarkdownDescription: "The Docker image tag for the container run on the Pod.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"compute_type": schema.StringAttribute{
				MarkdownDescription: "Set to GPU to create a GPU Pod. Set to CPU to create a CPU Pod.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("GPU"),
			},
			"cloud_type": schema.StringAttribute{
				MarkdownDescription: "Set to SECURE to create the Pod in Secure Cloud. Set to COMMUNITY to create the Pod in Community Cloud.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("SECURE"),
			},
			"gpu_count": schema.Int64Attribute{
				MarkdownDescription: "If the Pod is a GPU Pod, the number of GPUs attached to the Pod.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
			"vcpu_count": schema.Int64Attribute{
				MarkdownDescription: "If the Pod is a CPU Pod, the number of vCPUs allocated to the Pod.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(2),
			},
			"gpu_type_ids": schema.ListAttribute{
				MarkdownDescription: "If the Pod is a GPU Pod, a list of RunPod GPU types which can be attached to the Pod.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"cpu_flavor_ids": schema.ListAttribute{
				MarkdownDescription: "If the Pod is a CPU Pod, a list of RunPod CPU flavors which can be attached to the Pod.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"data_center_ids": schema.ListAttribute{
				MarkdownDescription: "A list of RunPod data center IDs where the Pod can be located.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"container_disk_in_gb": schema.Int64Attribute{
				MarkdownDescription: "The amount of disk space, in gigabytes (GB), to allocate on the container disk. Data is wiped when the Pod restarts.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(50),
			},
			"volume_in_gb": schema.Int64Attribute{
				MarkdownDescription: "The amount of disk space, in gigabytes (GB), to allocate on the Pod volume. Data is persisted across Pod restarts.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(20),
			},
			"volume_mount_path": schema.StringAttribute{
				MarkdownDescription: "The absolute path where the network volume will be mounted in the filesystem.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("/workspace"),
			},
			"ports": schema.ListAttribute{
				MarkdownDescription: "A list of ports exposed on the Pod. Each port is formatted as [port number]/[protocol].",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"env": schema.MapAttribute{
				MarkdownDescription: "Environment variables for the Pod.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"docker_entrypoint": schema.ListAttribute{
				MarkdownDescription: "If specified, overrides the ENTRYPOINT for the Docker image run on the Pod.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"docker_start_cmd": schema.ListAttribute{
				MarkdownDescription: "If specified, overrides the start CMD for the Docker image run on the Pod.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"template_id": schema.StringAttribute{
				MarkdownDescription: "If the Pod is created with a template, the unique string identifying that template.",
				Optional:            true,
			},
			"network_volume_id": schema.StringAttribute{
				MarkdownDescription: "The unique string identifying the network volume to attach to the Pod.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"interruptible": schema.BoolAttribute{
				MarkdownDescription: "Set to true to create an interruptible or spot Pod. Can be rented at a lower cost but can be stopped at any time.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"locked": schema.BoolAttribute{
				MarkdownDescription: "Set to true to lock a Pod. Locking a Pod disables stopping or resetting the Pod.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"min_vcpu_per_gpu": schema.Int64Attribute{
				MarkdownDescription: "If the Pod is a GPU Pod, the minimum number of virtual CPUs allocated to the Pod for each GPU.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(2),
			},
			"min_ram_per_gpu": schema.Int64Attribute{
				MarkdownDescription: "If the Pod is a GPU Pod, the minimum amount of RAM, in gigabytes (GB), allocated to the Pod for each GPU.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(8),
			},
			"min_download_mbps": schema.Float64Attribute{
				MarkdownDescription: "The minimum download speed, in megabits per second (Mbps), for the Pod.",
				Optional:            true,
			},
			"min_upload_mbps": schema.Float64Attribute{
				MarkdownDescription: "The minimum upload speed, in megabits per second (Mbps), for the Pod.",
				Optional:            true,
			},
			"min_disk_bandwidth_mbps": schema.Float64Attribute{
				MarkdownDescription: "The minimum disk bandwidth, in megabytes per second (MBps), for the Pod.",
				Optional:            true,
			},
			"support_public_ip": schema.BoolAttribute{
				MarkdownDescription: "If the Pod is on Community Cloud, set to true if you need the Pod to expose a public IP address.",
				Optional:            true,
			},
			"global_networking": schema.BoolAttribute{
				MarkdownDescription: "Set to true to enable global networking for the Pod.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"allowed_cuda_versions": schema.ListAttribute{
				MarkdownDescription: "If the Pod is a GPU Pod, a list of acceptable CUDA versions on the Pod.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"country_codes": schema.ListAttribute{
				MarkdownDescription: "A list of country codes where the Pod can be located.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"gpu_type_priority": schema.StringAttribute{
				MarkdownDescription: "If the Pod is a GPU Pod, set to availability to respond to current GPU type availability. Set to custom to always try to rent GPU types in the order specified.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("availability"),
			},
			"cpu_flavor_priority": schema.StringAttribute{
				MarkdownDescription: "If the Pod is a CPU Pod, set to availability to respond to current CPU flavor availability. Set to custom to always try to rent CPU flavors in the order specified.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("availability"),
			},
			"data_center_priority": schema.StringAttribute{
				MarkdownDescription: "Set to availability to respond to current machine availability. Set to custom to always try to rent machines from data centers in the order specified.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("availability"),
			},
			"container_registry_auth_id": schema.StringAttribute{
				MarkdownDescription: "Registry credentials ID.",
				Optional:            true,
			},
			// Computed fields
			"desired_status": schema.StringAttribute{
				MarkdownDescription: "The current expected status of the Pod.",
				Computed:            true,
			},
			"public_ip": schema.StringAttribute{
				MarkdownDescription: "The public IP address of the Pod.",
				Computed:            true,
			},
			"machine_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the host machine the Pod is running on.",
				Computed:            true,
			},
			"actual_data_center": schema.StringAttribute{
				MarkdownDescription: "The actual data center where the Pod was deployed.",
				Computed:            true,
			},
			"cost_per_hr": schema.Float64Attribute{
				MarkdownDescription: "The cost in RunPod credits per hour of running the Pod.",
				Computed:            true,
			},
			"adjusted_cost_per_hr": schema.Float64Attribute{
				MarkdownDescription: "The effective cost in RunPod credits per hour of running the Pod, adjusted by active Savings Plans.",
				Computed:            true,
			},
			"memory_in_gb": schema.Float64Attribute{
				MarkdownDescription: "The amount of RAM, in gigabytes (GB), attached to the Pod.",
				Computed:            true,
			},
			"last_started_at": schema.StringAttribute{
				MarkdownDescription: "The UTC timestamp when the Pod was last started.",
				Computed:            true,
			},
		},
	}
}

func (r *PodResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PodResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Pod")

	// Build create input
	input := &PodCreateInput{
		Name:                    data.Name.ValueString(),
		ImageName:               data.ImageName.ValueString(),
		ComputeType:             data.ComputeType.ValueString(),
		CloudType:               data.CloudType.ValueString(),
		VolumeMountPath:         data.VolumeMountPath.ValueString(),
		TemplateId:              data.TemplateId.ValueString(),
		NetworkVolumeId:         data.NetworkVolumeId.ValueString(),
		GPUTypePriority:         data.GPUTypePriority.ValueString(),
		CPUFlavorPriority:       data.CPUFlavorPriority.ValueString(),
		DataCenterPriority:      data.DataCenterPriority.ValueString(),
		ContainerRegistryAuthId: data.ContainerRegistryAuthId.ValueString(),
	}

	// Handle integer pointers
	if !data.GPUCount.IsNull() {
		gpuCount := int(data.GPUCount.ValueInt64())
		input.GPUCount = &gpuCount
	}
	if !data.VCPUCount.IsNull() {
		vcpuCount := int(data.VCPUCount.ValueInt64())
		input.VCPUCount = &vcpuCount
	}
	if !data.ContainerDiskInGb.IsNull() {
		diskSize := int(data.ContainerDiskInGb.ValueInt64())
		input.ContainerDiskInGb = &diskSize
	}
	if !data.VolumeInGb.IsNull() {
		volumeSize := int(data.VolumeInGb.ValueInt64())
		input.VolumeInGb = &volumeSize
	}
	if !data.MinVCPUPerGPU.IsNull() {
		minVCPU := int(data.MinVCPUPerGPU.ValueInt64())
		input.MinVCPUPerGPU = &minVCPU
	}
	if !data.MinRAMPerGPU.IsNull() {
		minRAM := int(data.MinRAMPerGPU.ValueInt64())
		input.MinRAMPerGPU = &minRAM
	}

	// Handle boolean pointers
	if !data.Interruptible.IsNull() {
		interruptible := data.Interruptible.ValueBool()
		input.Interruptible = &interruptible
	}
	if !data.Locked.IsNull() {
		locked := data.Locked.ValueBool()
		input.Locked = &locked
	}
	if !data.SupportPublicIp.IsNull() {
		supportPublicIp := data.SupportPublicIp.ValueBool()
		input.SupportPublicIp = &supportPublicIp
	}
	if !data.GlobalNetworking.IsNull() {
		globalNetworking := data.GlobalNetworking.ValueBool()
		input.GlobalNetworking = &globalNetworking
	}

	// Handle float pointers
	if !data.MinDownloadMbps.IsNull() {
		minDownload := data.MinDownloadMbps.ValueFloat64()
		input.MinDownloadMbps = &minDownload
	}
	if !data.MinUploadMbps.IsNull() {
		minUpload := data.MinUploadMbps.ValueFloat64()
		input.MinUploadMbps = &minUpload
	}
	if !data.MinDiskBandwidthMBps.IsNull() {
		minDiskBandwidth := data.MinDiskBandwidthMBps.ValueFloat64()
		input.MinDiskBandwidthMBps = &minDiskBandwidth
	}

	// Handle string lists
	if !data.GPUTypeIds.IsNull() {
		resp.Diagnostics.Append(data.GPUTypeIds.ElementsAs(ctx, &input.GPUTypeIds, false)...)
	}
	if !data.CPUFlavorIds.IsNull() {
		resp.Diagnostics.Append(data.CPUFlavorIds.ElementsAs(ctx, &input.CPUFlavorIds, false)...)
	}
	if !data.DataCenterIds.IsNull() {
		resp.Diagnostics.Append(data.DataCenterIds.ElementsAs(ctx, &input.DataCenterIds, false)...)
	}
	if !data.Ports.IsNull() {
		resp.Diagnostics.Append(data.Ports.ElementsAs(ctx, &input.Ports, false)...)
	}
	if !data.DockerEntrypoint.IsNull() {
		resp.Diagnostics.Append(data.DockerEntrypoint.ElementsAs(ctx, &input.DockerEntrypoint, false)...)
	}
	if !data.DockerStartCmd.IsNull() {
		resp.Diagnostics.Append(data.DockerStartCmd.ElementsAs(ctx, &input.DockerStartCmd, false)...)
	}
	if !data.AllowedCudaVersions.IsNull() {
		resp.Diagnostics.Append(data.AllowedCudaVersions.ElementsAs(ctx, &input.AllowedCudaVersions, false)...)
	}
	if !data.CountryCodes.IsNull() {
		resp.Diagnostics.Append(data.CountryCodes.ElementsAs(ctx, &input.CountryCodes, false)...)
	}

	// Handle map
	if !data.Env.IsNull() {
		resp.Diagnostics.Append(data.Env.ElementsAs(ctx, &input.Env, false)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	pod, err := r.client.CreatePod(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create pod, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Created Pod", map[string]interface{}{"id": pod.ID})

	// Update state with response
	r.updateStateFromPod(ctx, &data, pod)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PodResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading Pod", map[string]interface{}{"id": data.ID.ValueString()})

	pod, err := r.client.GetPod(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read pod, got error: %s", err))
		return
	}

	r.updateStateFromPod(ctx, &data, pod)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PodResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating Pod", map[string]interface{}{"id": data.ID.ValueString()})

	// Build update input
	input := &PodUpdateInput{
		Name:                    data.Name.ValueString(),
		ImageName:               data.ImageName.ValueString(),
		VolumeMountPath:         data.VolumeMountPath.ValueString(),
		ContainerRegistryAuthId: data.ContainerRegistryAuthId.ValueString(),
	}

	// Handle integer pointers
	if !data.ContainerDiskInGb.IsNull() {
		diskSize := int(data.ContainerDiskInGb.ValueInt64())
		input.ContainerDiskInGb = &diskSize
	}
	if !data.VolumeInGb.IsNull() {
		volumeSize := int(data.VolumeInGb.ValueInt64())
		input.VolumeInGb = &volumeSize
	}

	// Handle boolean pointers
	if !data.Locked.IsNull() {
		locked := data.Locked.ValueBool()
		input.Locked = &locked
	}
	if !data.GlobalNetworking.IsNull() {
		globalNetworking := data.GlobalNetworking.ValueBool()
		input.GlobalNetworking = &globalNetworking
	}

	// Handle string lists
	if !data.Ports.IsNull() {
		resp.Diagnostics.Append(data.Ports.ElementsAs(ctx, &input.Ports, false)...)
	}
	if !data.DockerEntrypoint.IsNull() {
		resp.Diagnostics.Append(data.DockerEntrypoint.ElementsAs(ctx, &input.DockerEntrypoint, false)...)
	}
	if !data.DockerStartCmd.IsNull() {
		resp.Diagnostics.Append(data.DockerStartCmd.ElementsAs(ctx, &input.DockerStartCmd, false)...)
	}

	// Handle map
	if !data.Env.IsNull() {
		resp.Diagnostics.Append(data.Env.ElementsAs(ctx, &input.Env, false)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	pod, err := r.client.UpdatePod(ctx, data.ID.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update pod, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Updated Pod", map[string]interface{}{"id": pod.ID})

	// Update state with response
	r.updateStateFromPod(ctx, &data, pod)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PodResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting Pod", map[string]interface{}{"id": data.ID.ValueString()})

	err := r.client.DeletePod(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete pod, got error: %s", err))
		return
	}

	// Wait for pod to be deleted (with timeout)
	timeout := time.After(5 * time.Minute)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			resp.Diagnostics.AddWarning("Timeout", "Timed out waiting for pod to be deleted")
			return
		case <-ticker.C:
			_, err := r.client.GetPod(ctx, data.ID.ValueString())
			if err != nil {
				// Pod is deleted
				tflog.Trace(ctx, "Deleted Pod", map[string]interface{}{"id": data.ID.ValueString()})
				return
			}
		}
	}
}

func (r *PodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// updateStateFromPod updates the Terraform state from a Pod API response
func (r *PodResource) updateStateFromPod(ctx context.Context, data *PodResourceModel, pod *Pod) {
	data.ID = types.StringValue(pod.ID)
	data.Name = types.StringValue(pod.Name)

	// Only update ImageName if API returns a non-empty value
	// RunPod API doesn't return image name in GET responses, so preserve the planned value
	if pod.ImageName != "" {
		data.ImageName = types.StringValue(pod.ImageName)
	}

	data.DesiredStatus = types.StringValue(pod.DesiredStatus)
	data.PublicIp = types.StringValue(pod.PublicIp)
	data.MachineId = types.StringValue(pod.MachineId)
	// TODO: Figure out how to get actual data center from RunPod API
	// For now, we'll leave this empty until we find a way to determine the region
	data.ActualDataCenter = types.StringValue("")
	data.CostPerHr = types.Float64Value(pod.CostPerHr)
	data.AdjustedCostPerHr = types.Float64Value(pod.AdjustedCostPerHr)
	data.MemoryInGb = types.Float64Value(pod.MemoryInGb)
	data.LastStartedAt = types.StringValue(pod.LastStartedAt)

	if pod.VolumeInGb > 0 {
		data.VolumeInGb = types.Int64Value(int64(pod.VolumeInGb))
	}
	if pod.ContainerDiskInGb > 0 {
		data.ContainerDiskInGb = types.Int64Value(int64(pod.ContainerDiskInGb))
	}
	if pod.VolumeMountPath != "" {
		data.VolumeMountPath = types.StringValue(pod.VolumeMountPath)
	}
}
