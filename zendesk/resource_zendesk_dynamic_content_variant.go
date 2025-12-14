package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
)

func resourceZendeskDynamicContentVariant() *schema.Resource {
	return &schema.Resource{
		Description: `This defines the variants of a dynamic_content. In order to delete this resource, need to remove the zendesk_dynamic_content resource and use terraform state rm for zendesk_dynamic_content_variant`,
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return createDynamicContentVariant(ctx, d, zd)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return readDynamicContentVariant(ctx, d, zd)
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return updateDynamicContentVariant(ctx, d, zd)
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return deleteDynamicContentVariant(ctx, d, zd)
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content": {
				Type:        schema.TypeString,
				Description: "Content of the dynamic content variant",
				Required:    true,
			},
			"locale_id": {
				Type:        schema.TypeInt,
				Description: "Locale Id for the dynamic content variant",
				Required:    true,
			},
			"dynamic_content_item_id": {
				Type:        schema.TypeInt,
				Description: "Reference for the dynamic content item",
				Required:    true,
			},
			"default": {
				Type:        schema.TypeBool,
				Description: "default resources cannot be deleted; only managed resources can be set as default a particular zendesk_dynamic_content",
				Required:    true,
			},
		},
	}
}

// ID        int64     `json:"id,omitempty"`
// 	URL       string    `json:"url,omitempty"`
// 	Content   string    `json:"content"`
// 	LocaleID  int64     `json:"locale_id"`
// 	Outdated  bool      `json:"outdated,omitempty"`
// 	Active    bool      `json:"active,omitempty"`
// 	Default   bool      `json:"default,omitempty"`
// 	CreatedAt time.Time `json:"created_at,omitempty"`
// 	UpdatedAt time.Time `json:"updated_at,omitempty"`

func marshalDynamicContentVariant(dc DynamicContentVariant, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":                     dc.URL,
		"content":                 dc.Content,
		"locale_id":               dc.LocaleID,
		"dynamic_content_item_id": dc.DynamicContentItemID,
		"default":                 dc.Default,
	}

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalDynamicContentVariant(d identifiableGetterSetter) (DynamicContentVariant, error) {
	dc := DynamicContentVariant{}

	if v := d.Id(); v != "" {
		dcv_id, dc_id, err := split_id(v)
		if err != nil {
			return dc, fmt.Errorf("could not parse dc ID")
		}
		dc.ID = dcv_id
		dc.DynamicContentItemID = dc_id
	}

	if v, ok := d.GetOk("url"); ok {
		dc.URL = v.(string)
	}

	if v, ok := d.GetOk("default"); ok {
		dc.Default = v.(bool)
	}

	if v, ok := d.GetOk("locale_id"); ok {
		dc.LocaleID = int64(v.(int))
	}
	fmt.Println("DCV Unmarshalling ->")
	fmt.Printf("dc.LocaleID %d", dc.LocaleID)

	if v, ok := d.GetOk("dynamic_content_item_id"); ok {
		// first request is without id
		if (dc.DynamicContentItemID) == 0 {
			dc.DynamicContentItemID = int64(v.(int))
		}
	}

	if v, ok := d.GetOk("content"); ok {
		dc.Content = v.(string)
	}

	dc.Default = true

	fmt.Println("Unmarshalling DCV")
	fmt.Printf("d.Id: %s", d.Id())
	fmt.Printf("dc.ID: %d", dc.ID)
	fmt.Printf("dc.DynamicContentItemID: %d", dc.DynamicContentItemID)
	return dc, nil
}

func createDynamicContentVariant(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	dc, err := unmarshalDynamicContentVariant(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API Request
	dc, err = Create(ctx, zd, dc)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d+%d", dc.ID, dc.DynamicContentItemID))

	err = marshalDynamicContentVariant(dc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readDynamicContentVariant(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	dcv_id, dc_id, err := split_id(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	dc, err := Get(ctx, zd, dc_id, dcv_id)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonData, err := json.Marshal(dc)
	if err != nil {
		return diag.FromErr(err)
	}
	fmt.Println("READ: MARSHALLED OUTPUT")
	fmt.Println(string(jsonData))

	err = marshalDynamicContentVariant(dc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateDynamicContentVariant(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	dc, err := unmarshalDynamicContentVariant(d)
	if err != nil {
		return diag.FromErr(err)
	}

	dcv_id, _, err := split_id(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API Request
	dc, err = Update(ctx, zd, dcv_id, dc)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalDynamicContentVariant(dc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func split_id(ids string) (int64, int64, error) {
	ids_arr := strings.Split(ids, "+")

	dcv_id, err := atoi64(ids_arr[0])
	if err != nil {
		return 0, 0, err
	}

	dc_id, err := atoi64(ids_arr[1])
	if err != nil {
		return 0, 0, err
	}

	return dcv_id, dc_id, err
}

func deleteDynamicContentVariant(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	dcv_id, dc_id, err := split_id(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = Delete(ctx, zd, dc_id, dcv_id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func Create(ctx context.Context, z *newClient.Client, field DynamicContentVariant) (DynamicContentVariant, error) {
	var result struct {
		DynamicContentVariant DynamicContentVariant `json:"variant"`
	}

	result.DynamicContentVariant = field

	fmt.Println("Create Payload: => ")
	fmt.Printf("%+v\n", result)
	fmt.Printf("%d dc_id\n", field.DynamicContentItemID)

	body, err := z.Post(ctx, fmt.Sprintf("/dynamic_content/items/%d/variants.json", field.DynamicContentItemID), result)

	if err != nil {
		return DynamicContentVariant{}, err
	}

	err = json.Unmarshal(body, &result)
	result.DynamicContentVariant.DynamicContentItemID = field.DynamicContentItemID
	if err != nil {
		return DynamicContentVariant{}, err
	}
	return result.DynamicContentVariant, nil
}

func Get(ctx context.Context, z *newClient.Client, dynamicContentItemID int64, viewID int64) (DynamicContentVariant, error) {
	var result struct {
		DynamicContentVariant DynamicContentVariant `json:"variant"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/dynamic_content/items/%d/variants/%d.json", dynamicContentItemID, viewID))
	fmt.Println("GET bar")
	fmt.Println(string(body))

	if err != nil {
		return DynamicContentVariant{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return DynamicContentVariant{}, err
	}

	result.DynamicContentVariant.DynamicContentItemID = dynamicContentItemID
	return result.DynamicContentVariant, err
}

// UpdateDynamicContentVariant updates a field with the specified ticket field
// ref: https://developer.zendesk.com/rest_api/docs/support/user_fields#update-ticket-field
func Update(ctx context.Context, z *newClient.Client, ticketID int64, field DynamicContentVariant) (DynamicContentVariant, error) {
	var result struct {
		DynamicContentVariant DynamicContentVariant `json:"variant"`
	}

	result.DynamicContentVariant = field

	fmt.Println("Update Processed payload: JSON")
	fmt.Printf("%+v\n", result)

	body, err := z.Put(ctx, fmt.Sprintf("/dynamic_content/items/%d/variants/%d.json", field.DynamicContentItemID, ticketID), result)
	fmt.Println("Update bar")
	fmt.Println(string(body))

	if err != nil {
		fmt.Println("Printing Error")
		fmt.Printf("%+v\n", err)
		return DynamicContentVariant{}, err
	}

	err = json.Unmarshal(body, &result)
	result.DynamicContentVariant.DynamicContentItemID = field.DynamicContentItemID
	if err != nil {
		return DynamicContentVariant{}, err
	}

	return result.DynamicContentVariant, err
}

// DeleteDynamicContentVariant deletes the specified ticket field
// ref: https://developer.zendesk.com/rest_api/docs/support/user_fields#Delete-ticket-field
func Delete(ctx context.Context, z *newClient.Client, dynamicContentItemID int64, viewID int64) error {
	fmt.Println("Deleting something")
	url := fmt.Sprintf("/dynamic_content/items/%d/variants/%d.json", dynamicContentItemID, viewID)
	fmt.Printf("delete at: %s", url)
	err := z.Delete(ctx, fmt.Sprintf("/dynamic_content/items/%d/variants/%d.json", dynamicContentItemID, viewID))

	if err != nil {
		return err
	}

	return nil
}

// https://developer.zendesk.com/rest_api/docs/support/dynamic_content#json-format-for-variants
type DynamicContentVariant struct {
	ID                   int64     `json:"id,omitempty"`
	URL                  string    `json:"url,omitempty"`
	Content              string    `json:"content"`
	LocaleID             int64     `json:"locale_id"`
	DynamicContentItemID int64     `json:",omitempty"` // needed to use this model because it didn't include this field
	Outdated             bool      `json:"outdated,omitempty"`
	Active               bool      `json:"active,omitempty"`
	Default              bool      `json:"default,omitempty"`
	CreatedAt            time.Time `json:"created_at,omitempty"`
	UpdatedAt            time.Time `json:"updated_at,omitempty"`
}
