package zendesk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
)

func dataSourceZendeskTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return readTagsDataSource(ctx, d, zd)
		},

		Schema: map[string]*schema.Schema{
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of tag names",
			},
			"count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total count of tags",
			},
		},
	}
}

func readTagsDataSource(ctx context.Context, d *schema.ResourceData, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	tags, err := zd.GetTags(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	count, err := zd.GetTagCount(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}

	d.SetId("tags")
	d.Set("tags", tagNames)
	d.Set("count", int(count))

	return diags
}

