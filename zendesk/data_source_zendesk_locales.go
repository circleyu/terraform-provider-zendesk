package zendesk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
)

func dataSourceZendeskLocales() *schema.Resource {
	return &schema.Resource{
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return readLocalesDataSource(ctx, d, zd)
		},

		Schema: map[string]*schema.Schema{
			"locales": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"locale": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Description: "List of available locales",
			},
		},
	}
}

func readLocalesDataSource(ctx context.Context, d *schema.ResourceData, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	locales, err := zd.GetLocales(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	localeList := make([]map[string]interface{}, len(locales))
	for i, locale := range locales {
		localeList[i] = map[string]interface{}{
			"id":         int(locale.ID),
			"url":        locale.URL,
			"locale":     locale.Locale,
			"name":       locale.Name,
			"created_at": locale.CreatedAt,
			"updated_at": locale.UpdatedAt,
		}
	}

	d.SetId("locales")
	d.Set("locales", localeList)

	return diags
}

