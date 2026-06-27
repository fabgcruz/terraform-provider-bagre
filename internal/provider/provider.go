package provider

import (
	"context"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure BagreProvider satisfies the provider.Provider interface.
var _ provider.Provider = &BagreProvider{}

type BagreProvider struct {
	version string
}

// BagreProviderModel maps the provider configuration block.
type BagreProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIToken types.String `tfsdk:"api_token"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &BagreProvider{version: version}
	}
}

func (p *BagreProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bagre"
	resp.Version = p.version
}

func (p *BagreProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage Bagre IPAM (IP Address Management) as code. Works with OpenTofu and Terraform.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "Base URL of the Bagre instance, e.g. https://ipam.example.com. May also be set via the BAGRE_ENDPOINT environment variable.",
			},
			"api_token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "API token (bagre_…) used as a Bearer credential. May also be set via the BAGRE_TOKEN environment variable.",
			},
		},
	}
}

func (p *BagreProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config BagreProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Config attribute wins; otherwise fall back to the environment variable.
	endpoint := os.Getenv("BAGRE_ENDPOINT")
	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}
	token := os.Getenv("BAGRE_TOKEN")
	if !config.APIToken.IsNull() {
		token = config.APIToken.ValueString()
	}

	endpoint = strings.TrimSpace(endpoint)
	token = strings.TrimSpace(token)

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing Bagre endpoint",
			"Set the provider `endpoint` attribute or the BAGRE_ENDPOINT environment variable.",
		)
	}
	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing Bagre API token",
			"Set the provider `api_token` attribute or the BAGRE_TOKEN environment variable. "+
				"Generate one in Bagre under Tokens de API (/admin/api-tokens).",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	client := NewClient(endpoint, token, p.version)
	resp.ResourceData = client
	resp.DataSourceData = client
}

func (p *BagreProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSiteResource,
	}
}

func (p *BagreProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
