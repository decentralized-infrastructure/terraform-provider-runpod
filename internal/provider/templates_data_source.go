package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TemplatesDataSource{}

func NewTemplatesDataSource() datasource.DataSource {
	return &TemplatesDataSource{}
}

type TemplatesDataSource struct {
	client *Client
}

type TemplatesDataSourceModel struct {
	Templates []TemplateDataModel `tfsdk:"templates"`
}

type TemplateDataModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	ImageName    types.String `tfsdk:"image_name"`
	IsServerless types.Bool   `tfsdk:"is_serverless"`
	IsPublic     types.Bool   `tfsdk:"is_public"`
	IsRunpod     types.Bool   `tfsdk:"is_runpod"`
}

func (d *TemplatesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_templates"
}

func (d *TemplatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source to list all RunPod Templates.",

		Attributes: map[string]schema.Attribute{
			"templates": schema.ListNestedAttribute{
				MarkdownDescription: "List of Templates.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier of the Template.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the Template.",
							Computed:            true,
						},
						"image_name": schema.StringAttribute{
							MarkdownDescription: "The Docker image for the Template.",
							Computed:            true,
						},
						"is_serverless": schema.BoolAttribute{
							MarkdownDescription: "Whether the template is for serverless endpoints.",
							Computed:            true,
						},
						"is_public": schema.BoolAttribute{
							MarkdownDescription: "Whether the template is public.",
							Computed:            true,
						},
						"is_runpod": schema.BoolAttribute{
							MarkdownDescription: "Whether the template is an official RunPod template.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *TemplatesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (d *TemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TemplatesDataSourceModel

	// Initialize empty slice
	data.Templates = []TemplateDataModel{}

	tflog.Debug(ctx, "Reading Templates data source")

	templates, err := d.client.ListTemplates(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list templates, got error: %s", err))
		return
	}

	for _, template := range templates {
		templateData := TemplateDataModel{
			ID:           types.StringValue(template.ID),
			Name:         types.StringValue(template.Name),
			ImageName:    types.StringValue(template.ImageName),
			IsServerless: types.BoolValue(template.IsServerless),
			IsPublic:     types.BoolValue(template.IsPublic),
			IsRunpod:     types.BoolValue(template.IsRunpod),
		}
		data.Templates = append(data.Templates, templateData)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
