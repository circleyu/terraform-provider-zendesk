package zendesk

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	client "github.com/nukosuke/go-zendesk/zendesk"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
)

func resourceZendeskDynamicContent() *schema.Resource {
	return &schema.Resource{
		Description: `Due to limitation of zendesk API creates a placeholder dynamic_content_variant with placeholder text`,
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return createDynamicContent(ctx, d, zd)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return readDynamicContent(ctx, d, zd)
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return updateDynamicContent(ctx, d, zd)
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return deleteDynamicContent(ctx, d, zd)
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Dynamic Content Item",
			},
			"locale_id": {
				Type:        schema.TypeInt,
				Description: "Default Locale Id for the dynamic content item",
				Computed:    true,
			},
		},
	}
}

func marshalDynamicContent(dc client.DynamicContentItem, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":  dc.URL,
		"name": dc.Name,
		// "content":   dc.Variants[0].Content,
		"locale_id": dc.DefaultLocaleID,
	}

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalDynamicContent(d identifiableGetterSetter) (client.DynamicContentItem, error) {
	dc := client.DynamicContentItem{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return dc, fmt.Errorf("could not parse dynamic content item id %s: %v", v, err)
		}
		dc.ID = id
	}

	if v, ok := d.GetOk("url"); ok {
		dc.URL = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		dc.Name = v.(string)
	}

	dc_variant := client.DynamicContentVariant{}
	dc_variant.Default = true // This lib only supports a single dc variant
	dc_variant.Content = "AUTO_GENERATED_CONTENT_ZENDESK_API_LIMITATION"
	// dc_variant.Active = false

	if v, ok := d.GetOk("locale_id"); ok {
		dc.DefaultLocaleID = int64(v.(int))
		dc_variant.LocaleID = dc.DefaultLocaleID
	} else {
		dc_variant.LocaleID = 16
	}
	dc.Variants = append(dc.Variants, dc_variant)

	return dc, nil
}

func createDynamicContent(ctx context.Context, d identifiableGetterSetter, zd client.DynamicContentAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	dc, err := unmarshalDynamicContent(d)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonData, err := json.Marshal(dc)
	fmt.Println("Create: => ")
	fmt.Println(string(jsonData))

	// Actual API Request
	dc, err = zd.CreateDynamicContentItem(ctx, dc)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", dc.ID))

	err = marshalDynamicContent(dc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readDynamicContent(ctx context.Context, d identifiableGetterSetter, zd client.DynamicContentAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	dc, err := zd.GetDynamicContentItem(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalDynamicContent(dc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateDynamicContent(ctx context.Context, d identifiableGetterSetter, zd client.DynamicContentAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	dc, err := unmarshalDynamicContent(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API Request
	dc, err = zd.UpdateDynamicContentItem(ctx, id, dc)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalDynamicContent(dc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func deleteDynamicContent(ctx context.Context, d identifiable, zd client.DynamicContentAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.DeleteDynamicContentItem(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
