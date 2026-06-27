package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &SiteResource{}
	_ resource.ResourceWithConfigure   = &SiteResource{}
	_ resource.ResourceWithImportState = &SiteResource{}
)

func NewSiteResource() resource.Resource {
	return &SiteResource{}
}

type SiteResource struct {
	client *Client
}

// SiteResourceModel maps a bagre_site resource block.
type SiteResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Code        types.String `tfsdk:"code"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (r *SiteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (r *SiteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A Bagre site (physical location / datacenter) that groups subnets.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Numeric identifier assigned by Bagre.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"code": schema.StringAttribute{
				Required:    true,
				Description: "Short unique code for the site, e.g. \"DC1\".",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Human-readable name.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Optional free-text description.",
			},
		},
	}
}

func (r *SiteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected provider data",
			fmt.Sprintf("expected *Client, got %T", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *SiteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SiteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.client.CreateSite(ctx, Site{
		Code:        plan.Code.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating site", err.Error())
		return
	}

	// Only the computed id is filled from the API; the rest already matches plan.
	plan.ID = types.StringValue(strconv.FormatInt(created.ID, 10))
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *SiteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SiteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid site id in state", err.Error())
		return
	}

	site, err := r.client.GetSite(ctx, id)
	if err != nil {
		if errors.Is(err, errNotFound) {
			resp.State.RemoveResource(ctx) // drifted/deleted out-of-band
			return
		}
		resp.Diagnostics.AddError("Error reading site", err.Error())
		return
	}

	state.Code = types.StringValue(site.Code)
	state.Name = types.StringValue(site.Name)
	// Keep description null if it was never set and the API has none, to avoid
	// null-vs-empty-string churn in the plan.
	if site.Description != "" {
		state.Description = types.StringValue(site.Description)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *SiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SiteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(plan.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid site id in state", err.Error())
		return
	}

	if _, err := r.client.UpdateSite(ctx, id, Site{
		Code:        plan.Code.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error updating site", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *SiteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SiteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid site id in state", err.Error())
		return
	}

	if err := r.client.DeleteSite(ctx, id); err != nil && !errors.Is(err, errNotFound) {
		resp.Diagnostics.AddError("Error deleting site", err.Error())
	}
}

func (r *SiteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
