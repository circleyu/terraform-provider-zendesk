package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
)

// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket-statuses/
func resourceZendeskCustomStatuses() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a custom status resource.",
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return createCustomStatus(ctx, d, zd)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return readCustomStatus(ctx, d, zd)
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return updateCustomStatus(ctx, d, zd)
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			// Custom statuses cannot be deleted via API
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Custom statuses cannot be deleted via API",
					Detail:   "The custom status will remain in Zendesk but will be removed from Terraform state.",
				},
			}
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "The API url of this custom status.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status_category": {
				Description: "The status category. Allowed values: new, open, pending, hold, solved.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"agent_label": {
				Description: "The label shown to agents.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"end_user_label": {
				Description: "The label shown to end users.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"description": {
				Description: "The description of the custom status.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"default": {
				Description: "Whether this is the default status for the category.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"active": {
				Description: "Whether the custom status is active.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"end_user_hidden": {
				Description: "Whether the status is hidden from end users.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"created_at": {
				Description: "The time the custom status was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "The time the custom status was last updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func marshalCustomStatus(status newClient.CustomStatus, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":             status.URL,
		"status_category": status.StatusCategory,
		"agent_label":     status.AgentLabel,
		"end_user_label":  status.EndUserLabel,
		"description":     status.Description,
		"default":         status.Default,
		"active":          status.Active,
		"end_user_hidden": status.EndUserHidden,
		"created_at":      status.CreatedAt,
		"updated_at":      status.UpdatedAt,
	}

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalCustomStatus(d identifiableGetterSetter) (newClient.CustomStatus, error) {
	status := newClient.CustomStatus{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return status, fmt.Errorf("could not parse custom status id %s: %v", v, err)
		}
		status.ID = id
	}

	if v, ok := d.GetOk("url"); ok {
		status.URL = v.(string)
	}

	if v, ok := d.GetOk("status_category"); ok {
		status.StatusCategory = v.(string)
	}

	if v, ok := d.GetOk("agent_label"); ok {
		status.AgentLabel = v.(string)
	}

	if v, ok := d.GetOk("end_user_label"); ok {
		status.EndUserLabel = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		status.Description = v.(string)
	}

	if v, ok := d.GetOk("default"); ok {
		status.Default = v.(bool)
	}

	if v, ok := d.GetOk("active"); ok {
		status.Active = v.(bool)
	}

	if v, ok := d.GetOk("end_user_hidden"); ok {
		status.EndUserHidden = v.(bool)
	}

	return status, nil
}

func createCustomStatus(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	status, err := unmarshalCustomStatus(d)
	if err != nil {
		return diag.FromErr(err)
	}

	status, err = zd.CreateCustomStatus(ctx, status)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", status.ID))

	err = marshalCustomStatus(status, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readCustomStatus(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	status, err := zd.GetCustomStatus(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalCustomStatus(status, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateCustomStatus(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	status, err := unmarshalCustomStatus(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	status, err = zd.UpdateCustomStatus(ctx, id, status)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalCustomStatus(status, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

