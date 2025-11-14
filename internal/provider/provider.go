package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &runpodProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &runpodProvider{
			version: version,
		}
	}
}

// runpodProvider is the provider implementation.
type runpodProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *runpodProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "runpod"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *runpodProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The RunPod API key. Can also be set via the RUNPOD_API_KEY environment variable.",
			},
		},
	}
}

// runpodProviderModel maps provider schema data to a Go type.
type runpodProviderModel struct {
	ApiKey types.String `tfsdk:"api_key"`
}

// Configure prepares a RunPod API client for data sources and resources.
func (p *runpodProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config runpodProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown RunPod API Key",
			"The provider cannot create the RunPod API client as there is an unknown configuration value for the RunPod API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the RUNPOD_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	api_key := os.Getenv("RUNPOD_API_KEY")

	if !config.ApiKey.IsNull() {
		api_key = config.ApiKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing RunPod API Key",
			"The provider cannot create the RunPod API client as there is a missing or empty value for the RunPod API key. "+
				"Set the API key value in the configuration or use the RUNPOD_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API client
	client := NewClient(api_key)
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *runpodProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPodsDataSource,
		NewEndpointsDataSource,
		NewNetworkVolumesDataSource,
		NewTemplatesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *runpodProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPodResource,
		NewEndpointResource,
		NewNetworkVolumeResource,
	}
}
