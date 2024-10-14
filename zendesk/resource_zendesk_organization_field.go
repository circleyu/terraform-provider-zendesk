package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/nukosuke/go-zendesk/zendesk"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
)

// https://developer.zendesk.com/rest_api/docs/core/organization_fields
func resourceZendeskOrganizationField() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a organization field resource.",
		CreateContext: resourceZendeskOrganizationFieldCreate,
		ReadContext:   resourceZendeskOrganizationFieldRead,
		UpdateContext: resourceZendeskOrganizationFieldUpdate,
		DeleteContext: resourceZendeskOrganizationFieldDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "The URL for this organization field.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "System or custom field type. Editable for custom field types and only on creation.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"dropdown", "lookup",
					"checkbox",
					"date",
					"decimal",
					"integer",
					"regexp",
					"text",
					"textarea",
				}, false),
			},
			"title": {
				Description: "The title of the organization field.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"key": {
				Description: "A unique key that identifies this custom field. This is used for updating the field and referencing in placeholders. The key must consist of only letters, numbers, and underscores. It can't be only numbers and can't be reused if deleted.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Describes the purpose of the organization field to users.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"position": {
				Description: "The relative position of the organization field on a ticket. Note that for accounts with ticket forms, positions are controlled by the different forms.",
				Type:        schema.TypeInt,
				Optional:    true,
				// positions 0 to 7 are reserved for system fields
				ValidateFunc: validation.IntAtLeast(8),
				Computed:     true,
			},
			"active": {
				Description: "Whether this field is available.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"regexp_for_validation": {
				Description: `For "regexp" fields only. The validation pattern for a field value to be deemed valid.`,
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"tag": {
				Description: `For "checkbox" fields only. A tag added to tickets when the checkbox field is selected.`,
				Type:        schema.TypeString,
				Optional:    true,
			},
			"custom_field_option": {
				Description: `Required and presented for a custom organization field of type "multiselect" or "tagger".`,
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Custom field option name.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"value": {
							Description: "Custom field option value.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "Custom field option id.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
				Optional: true,
			},
		},
	}
}

// marshalOrganizationField encodes the provided organization field into the provided resource data
func marshalOrganizationField(field client.OrganizationField, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":         field.URL,
		"type":        field.Type,
		"title":       field.Title,
		"key":         field.Key,
		"description": field.Description,
		"position":    field.Position,
		"active":      field.Active,
		// "required":              field.Required,
		// "collapsed_for_agents":  field.CollapsedForAgents,
		"regexp_for_validation": field.RegexpForValidation,
		// "title_in_portal":       field.TitleInPortal,
		// "visible_in_portal":     field.VisibleInPortal,
		// "editable_in_portal":    field.EditableInPortal,
		// "required_in_portal":    field.RequiredInPortal,
		"tag": field.Tag,
		// "sub_type_id":           field.SubTypeID,
		// "removable":             field.Removable,
		// "agent_description":     field.AgentDescription,
	}

	// set system field options
	// systemFieldOptions := make([]map[string]interface{}, 0)
	// for _, v := range field.SystemFieldOptions {
	// 	m := map[string]interface{}{
	// 		"name":  v.Name,
	// 		"value": v.Value,
	// 	}
	// 	systemFieldOptions = append(systemFieldOptions, m)
	// }

	// fields["system_field_options"] = systemFieldOptions

	// Set custom field options
	customFieldOptions := make([]map[string]interface{}, 0)
	for _, v := range field.CustomFieldOptions {
		m := map[string]interface{}{
			"name":  v.Name,
			"value": v.Value,
			"id":    v.ID,
		}
		customFieldOptions = append(customFieldOptions, m)
	}

	fields["custom_field_option"] = customFieldOptions

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

// unmarshalOrganizationField parses the provided ResourceData and returns a organization field
func unmarshalOrganizationField(d identifiableGetterSetter) (client.OrganizationField, error) {
	tf := client.OrganizationField{}

	if v := d.Id(); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return tf, fmt.Errorf("could not parse organization field id %s: %v", v, err)
		}
		tf.ID = id
	}

	if v, ok := d.GetOk("url"); ok {
		tf.URL = v.(string)
	}

	if v, ok := d.GetOk("type"); ok {
		tf.Type = v.(string)
	}

	if v, ok := d.GetOk("title"); ok {
		tf.Title = v.(string)
		tf.RawTitle = v.(string)
	}

	if v, ok := d.GetOk("key"); ok {
		tf.Key = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		tf.Description = v.(string)
		tf.RawDescription = v.(string)
	}

	if v, ok := d.GetOk("position"); ok {
		tf.Position = int64(v.(int))
	}

	if v, ok := d.GetOk("active"); ok {
		tf.Active = v.(bool)
	}

	// if v, ok := d.GetOk("required"); ok {
	// 	tf.Required = v.(bool)
	// }

	if v, ok := d.GetOk("regexp_for_validation"); ok {
		tf.RegexpForValidation = v.(string)
	}

	// if v, ok := d.GetOk("title_in_portal"); ok {
	// 	tf.TitleInPortal = v.(string)
	// 	tf.RawTitleInPortal = v.(string)
	// }

	// if v, ok := d.GetOk("visible_in_portal"); ok {
	// 	tf.VisibleInPortal = v.(bool)
	// }

	// if v, ok := d.GetOk("editable_in_portal"); ok {
	// 	tf.EditableInPortal = v.(bool)
	// }

	// if v, ok := d.GetOk("required_in_portal"); ok {
	// 	tf.RequiredInPortal = v.(bool)
	// }

	if v, ok := d.GetOk("tag"); ok {
		tf.Tag = v.(string)
	}

	// if v, ok := d.GetOk("sub_type_id"); ok {
	// 	tf.SubTypeID = int64(v.(int))
	// }

	// if v, ok := d.GetOk("removable"); ok {
	// 	tf.Removable = v.(bool)
	// }

	// if v, ok := d.GetOk("agent_description"); ok {
	// 	tf.AgentDescription = v.(string)
	// }

	if v, ok := d.GetOk("custom_field_option"); ok {
		options := v.(*schema.Set).List()
		customFieldOptions := make([]client.CustomFieldOption, 0)
		for _, o := range options {
			option, ok := o.(map[string]interface{})
			if !ok {
				return tf, fmt.Errorf("could not parse custom options for field %v", tf)
			}

			customFieldOptions = append(customFieldOptions, client.CustomFieldOption{
				Name:  option["name"].(string),
				Value: option["value"].(string),
				ID:    int64(option["id"].(int)),
			})
		}

		tf.CustomFieldOptions = customFieldOptions
	}

	// if v, ok := d.GetOk("system_field_options"); ok {
	// 	options := v.(*schema.Set).List()
	// 	systemFieldOptions := make([]client.OrganizationFieldSystemFieldOption, 0)
	// 	for _, o := range options {
	// 		option, ok := o.(map[string]interface{})
	// 		if !ok {
	// 			return tf, fmt.Errorf("could not parse system options for field %v", tf)
	// 		}

	// 		systemFieldOptions = append(systemFieldOptions, client.OrganizationFieldSystemFieldOption{
	// 			Name:  option["name"].(string),
	// 			Value: option["value"].(string),
	// 		})
	// 	}

	// 	tf.SystemFieldOptions = systemFieldOptions
	// }

	return tf, nil
}

func resourceZendeskOrganizationFieldCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return createOrganizationField(ctx, d, zd)
}

func createOrganizationField(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	tf, err := unmarshalOrganizationField(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API request
	tf, err = zd.CreateOrganizationField(ctx, tf)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", tf.ID))

	err = marshalOrganizationField(tf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZendeskOrganizationFieldRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return readOrganizationField(ctx, d, zd)
}

func readOrganizationField(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	field, err := GetOrganizationField(ctx, zd, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalOrganizationField(field, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZendeskOrganizationFieldUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return updateOrganizationField(ctx, d, zd)
}

func updateOrganizationField(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	tf, err := unmarshalOrganizationField(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API request
	tf, err = UpdateOrganizationField(ctx, zd, id, tf)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalOrganizationField(tf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZendeskOrganizationFieldDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return deleteOrganizationField(ctx, d, zd)
}

func deleteOrganizationField(ctx context.Context, d identifiable, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	err = DeleteOrganizationField(ctx, zd, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// GetOrganizationField gets a specified ticket field
// ref: https://developer.zendesk.com/rest_api/docs/support/organization_fields#show-ticket-field
func GetOrganizationField(ctx context.Context, z *newClient.Client, organizationID int64) (client.OrganizationField, error) {
	var result struct {
		OrganizationField client.OrganizationField `json:"organization_field"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/organization_fields/%d.json", organizationID))

	if err != nil {
		return client.OrganizationField{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return client.OrganizationField{}, err
	}

	return result.OrganizationField, err
}

// UpdateOrganizationField updates a field with the specified ticket field
// ref: https://developer.zendesk.com/rest_api/docs/support/organization_fields#update-ticket-field
func UpdateOrganizationField(ctx context.Context, z *newClient.Client, ticketID int64, field client.OrganizationField) (client.OrganizationField, error) {
	var result, data struct {
		OrganizationField client.OrganizationField `json:"organization_field"`
	}

	data.OrganizationField = field

	body, err := z.Put(ctx, fmt.Sprintf("/organization_fields/%d.json", ticketID), data)

	if err != nil {
		return client.OrganizationField{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return client.OrganizationField{}, err
	}

	return result.OrganizationField, err
}

// DeleteOrganizationField deletes the specified ticket field
// ref: https://developer.zendesk.com/rest_api/docs/support/organization_fields#Delete-ticket-field
func DeleteOrganizationField(ctx context.Context, z *newClient.Client, ticketID int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/organization_fields/%d.json", ticketID))

	if err != nil {
		return err
	}

	return nil
}
