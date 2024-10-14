package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
)

// https://developer.zendesk.com/rest_api/docs/core/ticket_fields
func resourceZendeskTriggerCategory() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a trigger category resource.",
		CreateContext: resourceZendeskTriggerCategoryCreate,
		ReadContext:   resourceZendeskTriggerCategoryRead,
		UpdateContext: resourceZendeskTriggerCategoryUpdate,
		DeleteContext: resourceZendeskTriggerCategoryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the category",
				Type:        schema.TypeString,
				Required:    true,
			},
			"position": {
				Description:  "The relative position of the trigger category",
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
				Computed:     true,
			},
		},
	}
}

// marshalTriggerCategory encodes the provided ticket field into the provided resource data
func marshalTriggerCategory(triggerCategory TriggerCategory, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"name":     triggerCategory.Name,
		"position": triggerCategory.Position,
	}

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

// unmarshalTriggerCategory parses the provided ResourceData and returns a TriggerCategory
func unmarshalTriggerCategory(d identifiableGetterSetter) (TriggerCategory, error) {
	tf := TriggerCategory{}

	if v := d.Id(); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return tf, fmt.Errorf("could not parse trigger category id %s: %v", v, err)
		}
		tf.ID = id
	}

	if v, ok := d.GetOk("position"); ok {
		tf.Position = int64(v.(int))
	}

	if v, ok := d.GetOk("name"); ok {
		tf.Name = v.(string)
	}

	return tf, nil
}

func resourceZendeskTriggerCategoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return createTriggerCategory(ctx, d, zd)
}

func createTriggerCategory(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	fmt.Println("Creating Trigger Category")
	tf, err := unmarshalTriggerCategory(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API request
	tf, err = CreateTriggerCategory(ctx, zd, tf)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", tf.ID))

	err = marshalTriggerCategory(tf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZendeskTriggerCategoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return readTriggerCategory(ctx, d, zd)
}

func readTriggerCategory(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	triggerCategory, err := GetTriggerCategory(ctx, zd, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalTriggerCategory(triggerCategory, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZendeskTriggerCategoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return updateTriggerCategory(ctx, d, zd)
}

func updateTriggerCategory(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	tf, err := unmarshalTriggerCategory(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API request
	tf, err = UpdateTriggerCategory(ctx, zd, id, tf)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalTriggerCategory(tf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZendeskTriggerCategoryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return deleteTriggerCategory(ctx, d, zd)
}

func deleteTriggerCategory(ctx context.Context, d identifiable, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	err = DeleteTriggerCategory(ctx, zd, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

type (
	// TriggerCategory has a certain structure in Get & Different structure in
	// Put/Post
	TriggerCategory struct {
		ID       int64  `json:"id,string,omitempty"`
		Position int64  `json:"position"`
		Name     string `json:"name"`
	}
)

func CreateTriggerCategory(ctx context.Context, z *newClient.Client, triggerCategory TriggerCategory) (TriggerCategory, error) {
	var data, result struct {
		TriggerCategory TriggerCategory `json:"trigger_category"`
	}

	data.TriggerCategory = triggerCategory

	body, err := z.Post(ctx, "/trigger_categories.json", data)

	if err != nil {
		return TriggerCategory{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return TriggerCategory{}, err
	}
	return result.TriggerCategory, nil
}

func GetTriggerCategory(ctx context.Context, z *newClient.Client, TriggerCategoryID int64) (TriggerCategory, error) {
	var result struct {
		TriggerCategory TriggerCategory `json:"trigger_category"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/trigger_categories/%d.json", TriggerCategoryID))

	if err != nil {
		return TriggerCategory{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return TriggerCategory{}, err
	}

	return result.TriggerCategory, err
}

// UpdateTriggerCategory updates a field with the specified ticket field
// ref: https://developer.zendesk.com/rest_api/docs/support/user_fields#update-ticket-field
func UpdateTriggerCategory(ctx context.Context, z *newClient.Client, ticketID int64, triggerCategory TriggerCategory) (TriggerCategory, error) {
	var data, result struct {
		TriggerCategory TriggerCategory `json:"trigger_category"`
	}

	data.TriggerCategory = triggerCategory

	body, err := z.Put(ctx, fmt.Sprintf("/trigger_categories/%d.json", ticketID), data)

	if err != nil {
		return TriggerCategory{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return TriggerCategory{}, err
	}

	return result.TriggerCategory, err
}

func DeleteTriggerCategory(ctx context.Context, z *newClient.Client, TriggerCategoryID int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/trigger_categories/%d.json", TriggerCategoryID))

	if err != nil {
		return err
	}

	return nil
}
