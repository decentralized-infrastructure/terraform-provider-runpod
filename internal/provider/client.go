package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://rest.runpod.io/v1"
)

// Client is the RunPod API client
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new RunPod API client
func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL: defaultBaseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Minute * 5,
		},
	}
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %w", err)
	}

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

// Pod represents a RunPod Pod
type Pod struct {
	ID                      string                 `json:"id,omitempty"`
	Name                    string                 `json:"name,omitempty"`
	ImageName               string                 `json:"image,omitempty"`
	ComputeType             string                 `json:"computeType,omitempty"`
	CloudType               string                 `json:"cloudType,omitempty"`
	GPUCount                int                    `json:"gpuCount,omitempty"`
	VCPUCount               int                    `json:"vcpuCount,omitempty"`
	MemoryInGb              float64                `json:"memoryInGb,omitempty"`
	GPUTypeIds              []string               `json:"gpuTypeIds,omitempty"`
	CPUFlavorIds            []string               `json:"cpuFlavorIds,omitempty"`
	DataCenterIds           []string               `json:"dataCenterIds,omitempty"`
	ContainerDiskInGb       int                    `json:"containerDiskInGb,omitempty"`
	VolumeInGb              int                    `json:"volumeInGb,omitempty"`
	VolumeMountPath         string                 `json:"volumeMountPath,omitempty"`
	Ports                   []string               `json:"ports,omitempty"`
	Env                     map[string]string      `json:"env,omitempty"`
	DockerEntrypoint        []string               `json:"dockerEntrypoint,omitempty"`
	DockerStartCmd          []string               `json:"dockerStartCmd,omitempty"`
	TemplateId              string                 `json:"templateId,omitempty"`
	NetworkVolumeId         string                 `json:"networkVolumeId,omitempty"`
	Interruptible           bool                   `json:"interruptible,omitempty"`
	Locked                  bool                   `json:"locked,omitempty"`
	MinVCPUPerGPU           int                    `json:"minVCPUPerGPU,omitempty"`
	MinRAMPerGPU            int                    `json:"minRAMPerGPU,omitempty"`
	MinDownloadMbps         float64                `json:"minDownloadMbps,omitempty"`
	MinUploadMbps           float64                `json:"minUploadMbps,omitempty"`
	MinDiskBandwidthMBps    float64                `json:"minDiskBandwidthMBps,omitempty"`
	SupportPublicIp         *bool                  `json:"supportPublicIp,omitempty"`
	GlobalNetworking        bool                   `json:"globalNetworking,omitempty"`
	AllowedCudaVersions     []string               `json:"allowedCudaVersions,omitempty"`
	CountryCodes            []string               `json:"countryCodes,omitempty"`
	GPUTypePriority         string                 `json:"gpuTypePriority,omitempty"`
	CPUFlavorPriority       string                 `json:"cpuFlavorPriority,omitempty"`
	DataCenterPriority      string                 `json:"dataCenterPriority,omitempty"`
	ContainerRegistryAuthId string                 `json:"containerRegistryAuthId,omitempty"`
	DesiredStatus           string                 `json:"desiredStatus,omitempty"`
	PublicIp                string                 `json:"publicIp,omitempty"`
	PortMappings            map[string]int         `json:"portMappings,omitempty"`
	MachineId               string                 `json:"machineId,omitempty"`
	CostPerHr               float64                `json:"costPerHr,omitempty"`
	AdjustedCostPerHr       float64                `json:"adjustedCostPerHr,omitempty"`
	LastStartedAt           string                 `json:"lastStartedAt,omitempty"`
	LastStatusChange        string                 `json:"lastStatusChange,omitempty"`
	Machine                 map[string]interface{} `json:"machine,omitempty"`
	GPU                     map[string]interface{} `json:"gpu,omitempty"`
	NetworkVolume           map[string]interface{} `json:"networkVolume,omitempty"`
}

// PodCreateInput represents the input for creating a Pod
type PodCreateInput struct {
	Name                    string            `json:"name,omitempty"`
	ImageName               string            `json:"imageName,omitempty"`
	ComputeType             string            `json:"computeType,omitempty"`
	CloudType               string            `json:"cloudType,omitempty"`
	GPUCount                *int              `json:"gpuCount,omitempty"`
	VCPUCount               *int              `json:"vcpuCount,omitempty"`
	GPUTypeIds              []string          `json:"gpuTypeIds,omitempty"`
	CPUFlavorIds            []string          `json:"cpuFlavorIds,omitempty"`
	DataCenterIds           []string          `json:"dataCenterIds,omitempty"`
	ContainerDiskInGb       *int              `json:"containerDiskInGb,omitempty"`
	VolumeInGb              *int              `json:"volumeInGb,omitempty"`
	VolumeMountPath         string            `json:"volumeMountPath,omitempty"`
	Ports                   []string          `json:"ports,omitempty"`
	Env                     map[string]string `json:"env,omitempty"`
	DockerEntrypoint        []string          `json:"dockerEntrypoint,omitempty"`
	DockerStartCmd          []string          `json:"dockerStartCmd,omitempty"`
	TemplateId              string            `json:"templateId,omitempty"`
	NetworkVolumeId         string            `json:"networkVolumeId,omitempty"`
	Interruptible           *bool             `json:"interruptible,omitempty"`
	Locked                  *bool             `json:"locked,omitempty"`
	MinVCPUPerGPU           *int              `json:"minVCPUPerGPU,omitempty"`
	MinRAMPerGPU            *int              `json:"minRAMPerGPU,omitempty"`
	MinDownloadMbps         *float64          `json:"minDownloadMbps,omitempty"`
	MinUploadMbps           *float64          `json:"minUploadMbps,omitempty"`
	MinDiskBandwidthMBps    *float64          `json:"minDiskBandwidthMBps,omitempty"`
	SupportPublicIp         *bool             `json:"supportPublicIp,omitempty"`
	GlobalNetworking        *bool             `json:"globalNetworking,omitempty"`
	AllowedCudaVersions     []string          `json:"allowedCudaVersions,omitempty"`
	CountryCodes            []string          `json:"countryCodes,omitempty"`
	GPUTypePriority         string            `json:"gpuTypePriority,omitempty"`
	CPUFlavorPriority       string            `json:"cpuFlavorPriority,omitempty"`
	DataCenterPriority      string            `json:"dataCenterPriority,omitempty"`
	ContainerRegistryAuthId string            `json:"containerRegistryAuthId,omitempty"`
}

// PodUpdateInput represents the input for updating a Pod (triggers reset)
type PodUpdateInput struct {
	Name                    string            `json:"name,omitempty"`
	ImageName               string            `json:"imageName,omitempty"`
	ContainerDiskInGb       *int              `json:"containerDiskInGb,omitempty"`
	VolumeInGb              *int              `json:"volumeInGb,omitempty"`
	VolumeMountPath         string            `json:"volumeMountPath,omitempty"`
	Ports                   []string          `json:"ports,omitempty"`
	Env                     map[string]string `json:"env,omitempty"`
	DockerEntrypoint        []string          `json:"dockerEntrypoint,omitempty"`
	DockerStartCmd          []string          `json:"dockerStartCmd,omitempty"`
	Locked                  *bool             `json:"locked,omitempty"`
	GlobalNetworking        *bool             `json:"globalNetworking,omitempty"`
	ContainerRegistryAuthId string            `json:"containerRegistryAuthId,omitempty"`
}

// PodUpdateInPlaceInput represents the input for updating a Pod in place (no reset)
type PodUpdateInPlaceInput struct {
	Name   string `json:"name,omitempty"`
	Locked *bool  `json:"locked,omitempty"`
}

// CreatePod creates a new Pod
func (c *Client) CreatePod(ctx context.Context, input *PodCreateInput) (*Pod, error) {
	resp, err := c.doRequest(ctx, "POST", "/pods", input)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pod Pod
	if err := json.NewDecoder(resp.Body).Decode(&pod); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &pod, nil
}

// GetPod retrieves a Pod by ID
func (c *Client) GetPod(ctx context.Context, id string) (*Pod, error) {
	resp, err := c.doRequest(ctx, "GET", "/pods/"+id, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pod Pod
	if err := json.NewDecoder(resp.Body).Decode(&pod); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &pod, nil
}

// UpdatePod updates a Pod (triggers reset)
func (c *Client) UpdatePod(ctx context.Context, id string, input *PodUpdateInput) (*Pod, error) {
	resp, err := c.doRequest(ctx, "PUT", "/pods/"+id, input)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pod Pod
	if err := json.NewDecoder(resp.Body).Decode(&pod); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &pod, nil
}

// UpdatePodInPlace updates a Pod without triggering a reset
func (c *Client) UpdatePodInPlace(ctx context.Context, id string, input *PodUpdateInPlaceInput) (*Pod, error) {
	resp, err := c.doRequest(ctx, "PATCH", "/pods/"+id, input)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pod Pod
	if err := json.NewDecoder(resp.Body).Decode(&pod); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &pod, nil
}

// DeletePod terminates a Pod
func (c *Client) DeletePod(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, "DELETE", "/pods/"+id, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// StopPod stops a Pod
func (c *Client) StopPod(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, "POST", "/pods/"+id+"/stop", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// StartPod starts a Pod
func (c *Client) StartPod(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, "POST", "/pods/"+id+"/start", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Endpoint represents a RunPod Serverless Endpoint
type Endpoint struct {
	ID                  string                 `json:"id,omitempty"`
	Name                string                 `json:"name,omitempty"`
	TemplateId          string                 `json:"templateId,omitempty"`
	ComputeType         string                 `json:"computeType,omitempty"`
	GPUCount            int                    `json:"gpuCount,omitempty"`
	VCPUCount           int                    `json:"vcpuCount,omitempty"`
	GPUTypeIds          []string               `json:"gpuTypeIds,omitempty"`
	CPUFlavorIds        []string               `json:"cpuFlavorIds,omitempty"`
	DataCenterIds       []string               `json:"dataCenterIds,omitempty"`
	NetworkVolumeId     string                 `json:"networkVolumeId,omitempty"`
	WorkersMin          int                    `json:"workersMin,omitempty"`
	WorkersMax          int                    `json:"workersMax,omitempty"`
	IdleTimeout         int                    `json:"idleTimeout,omitempty"`
	ExecutionTimeoutMs  int                    `json:"executionTimeoutMs,omitempty"`
	ScalerType          string                 `json:"scalerType,omitempty"`
	ScalerValue         int                    `json:"scalerValue,omitempty"`
	AllowedCudaVersions []string               `json:"allowedCudaVersions,omitempty"`
	Env                 map[string]string      `json:"env,omitempty"`
	Flashboot           bool                   `json:"flashboot,omitempty"`
	CreatedAt           string                 `json:"createdAt,omitempty"`
	UserId              string                 `json:"userId,omitempty"`
	Version             int                    `json:"version,omitempty"`
	Template            map[string]interface{} `json:"template,omitempty"`
}

// EndpointCreateInput represents the input for creating an Endpoint
type EndpointCreateInput struct {
	Name                string   `json:"name,omitempty"`
	TemplateId          string   `json:"templateId"`
	ComputeType         string   `json:"computeType,omitempty"`
	GPUCount            *int     `json:"gpuCount,omitempty"`
	VCPUCount           *int     `json:"vcpuCount,omitempty"`
	GPUTypeIds          []string `json:"gpuTypeIds,omitempty"`
	CPUFlavorIds        []string `json:"cpuFlavorIds,omitempty"`
	DataCenterIds       []string `json:"dataCenterIds,omitempty"`
	NetworkVolumeId     string   `json:"networkVolumeId,omitempty"`
	WorkersMin          *int     `json:"workersMin,omitempty"`
	WorkersMax          *int     `json:"workersMax,omitempty"`
	IdleTimeout         *int     `json:"idleTimeout,omitempty"`
	ExecutionTimeoutMs  *int     `json:"executionTimeoutMs,omitempty"`
	ScalerType          string   `json:"scalerType,omitempty"`
	ScalerValue         *int     `json:"scalerValue,omitempty"`
	AllowedCudaVersions []string `json:"allowedCudaVersions,omitempty"`
	Flashboot           *bool    `json:"flashboot,omitempty"`
}

// EndpointUpdateInput represents the input for updating an Endpoint
type EndpointUpdateInput struct {
	Name                string   `json:"name,omitempty"`
	TemplateId          string   `json:"templateId,omitempty"`
	GPUCount            *int     `json:"gpuCount,omitempty"`
	VCPUCount           *int     `json:"vcpuCount,omitempty"`
	GPUTypeIds          []string `json:"gpuTypeIds,omitempty"`
	CPUFlavorIds        []string `json:"cpuFlavorIds,omitempty"`
	DataCenterIds       []string `json:"dataCenterIds,omitempty"`
	NetworkVolumeId     string   `json:"networkVolumeId,omitempty"`
	WorkersMin          *int     `json:"workersMin,omitempty"`
	WorkersMax          *int     `json:"workersMax,omitempty"`
	IdleTimeout         *int     `json:"idleTimeout,omitempty"`
	ExecutionTimeoutMs  *int     `json:"executionTimeoutMs,omitempty"`
	ScalerType          string   `json:"scalerType,omitempty"`
	ScalerValue         *int     `json:"scalerValue,omitempty"`
	AllowedCudaVersions []string `json:"allowedCudaVersions,omitempty"`
	Flashboot           *bool    `json:"flashboot,omitempty"`
}

// CreateEndpoint creates a new Endpoint
func (c *Client) CreateEndpoint(ctx context.Context, input *EndpointCreateInput) (*Endpoint, error) {
	resp, err := c.doRequest(ctx, "POST", "/endpoints", input)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var endpoint Endpoint
	if err := json.NewDecoder(resp.Body).Decode(&endpoint); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &endpoint, nil
}

// GetEndpoint retrieves an Endpoint by ID
func (c *Client) GetEndpoint(ctx context.Context, id string) (*Endpoint, error) {
	resp, err := c.doRequest(ctx, "GET", "/endpoints/"+id, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var endpoint Endpoint
	if err := json.NewDecoder(resp.Body).Decode(&endpoint); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &endpoint, nil
}

// UpdateEndpoint updates an Endpoint
func (c *Client) UpdateEndpoint(ctx context.Context, id string, input *EndpointUpdateInput) (*Endpoint, error) {
	resp, err := c.doRequest(ctx, "PATCH", "/endpoints/"+id, input)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var endpoint Endpoint
	if err := json.NewDecoder(resp.Body).Decode(&endpoint); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &endpoint, nil
}

// DeleteEndpoint deletes an Endpoint
func (c *Client) DeleteEndpoint(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, "DELETE", "/endpoints/"+id, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// ListEndpoints lists all Endpoints
func (c *Client) ListEndpoints(ctx context.Context) ([]Endpoint, error) {
	resp, err := c.doRequest(ctx, "GET", "/endpoints", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var endpoints []Endpoint
	if err := json.NewDecoder(resp.Body).Decode(&endpoints); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return endpoints, nil
}

// NetworkVolume represents a RunPod Network Volume
type NetworkVolume struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Size         int    `json:"size,omitempty"`
	DataCenterId string `json:"dataCenterId,omitempty"`
}

// NetworkVolumeCreateInput represents the input for creating a Network Volume
type NetworkVolumeCreateInput struct {
	Name         string `json:"name"`
	Size         int    `json:"size"`
	DataCenterId string `json:"dataCenterId"`
}

// NetworkVolumeUpdateInput represents the input for updating a Network Volume
type NetworkVolumeUpdateInput struct {
	Name string `json:"name,omitempty"`
	Size *int   `json:"size,omitempty"`
}

// CreateNetworkVolume creates a new Network Volume
func (c *Client) CreateNetworkVolume(ctx context.Context, input *NetworkVolumeCreateInput) (*NetworkVolume, error) {
	resp, err := c.doRequest(ctx, "POST", "/networkvolumes", input)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var volume NetworkVolume
	if err := json.NewDecoder(resp.Body).Decode(&volume); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &volume, nil
}

// GetNetworkVolume retrieves a Network Volume by ID
func (c *Client) GetNetworkVolume(ctx context.Context, id string) (*NetworkVolume, error) {
	resp, err := c.doRequest(ctx, "GET", "/networkvolumes/"+id, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var volume NetworkVolume
	if err := json.NewDecoder(resp.Body).Decode(&volume); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &volume, nil
}

// UpdateNetworkVolume updates a Network Volume
func (c *Client) UpdateNetworkVolume(ctx context.Context, id string, input *NetworkVolumeUpdateInput) (*NetworkVolume, error) {
	resp, err := c.doRequest(ctx, "PATCH", "/networkvolumes/"+id, input)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var volume NetworkVolume
	if err := json.NewDecoder(resp.Body).Decode(&volume); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &volume, nil
}

// DeleteNetworkVolume deletes a Network Volume
func (c *Client) DeleteNetworkVolume(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, "DELETE", "/networkvolumes/"+id, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// ListNetworkVolumes lists all Network Volumes
func (c *Client) ListNetworkVolumes(ctx context.Context) ([]NetworkVolume, error) {
	resp, err := c.doRequest(ctx, "GET", "/networkvolumes", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var volumes []NetworkVolume
	if err := json.NewDecoder(resp.Body).Decode(&volumes); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return volumes, nil
}

// ListPods lists all Pods
func (c *Client) ListPods(ctx context.Context) ([]Pod, error) {
	resp, err := c.doRequest(ctx, "GET", "/pods", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pods []Pod
	if err := json.NewDecoder(resp.Body).Decode(&pods); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return pods, nil
}

// Template represents a RunPod Template
type Template struct {
	ID                      string            `json:"id,omitempty"`
	Name                    string            `json:"name,omitempty"`
	ImageName               string            `json:"imageName,omitempty"`
	Category                string            `json:"category,omitempty"`
	ContainerDiskInGb       int               `json:"containerDiskInGb,omitempty"`
	VolumeInGb              int               `json:"volumeInGb,omitempty"`
	VolumeMountPath         string            `json:"volumeMountPath,omitempty"`
	Ports                   []string          `json:"ports,omitempty"`
	Env                     map[string]string `json:"env,omitempty"`
	DockerEntrypoint        []string          `json:"dockerEntrypoint,omitempty"`
	DockerStartCmd          []string          `json:"dockerStartCmd,omitempty"`
	IsPublic                bool              `json:"isPublic,omitempty"`
	IsRunpod                bool              `json:"isRunpod,omitempty"`
	IsServerless            bool              `json:"isServerless,omitempty"`
	Readme                  string            `json:"readme,omitempty"`
	RuntimeInMin            int               `json:"runtimeInMin,omitempty"`
	StartJupyter            bool              `json:"startJupyter,omitempty"`
	StartSsh                bool              `json:"startSsh,omitempty"`
	ContainerRegistryAuthId string            `json:"containerRegistryAuthId,omitempty"`
	Earned                  float64           `json:"earned,omitempty"`
}

// ListTemplates lists all Templates
func (c *Client) ListTemplates(ctx context.Context) ([]Template, error) {
	resp, err := c.doRequest(ctx, "GET", "/templates", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var templates []Template
	if err := json.NewDecoder(resp.Body).Decode(&templates); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return templates, nil
}
