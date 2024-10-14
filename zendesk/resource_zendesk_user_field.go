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

// https://developer.zendesk.com/rest_api/docs/core/user_fields
func resourceZendeskUserField() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a user field resource.",
		CreateContext: resourceZendeskUserFieldCreate,
		ReadContext:   resourceZendeskUserFieldRead,
		UpdateContext: resourceZendeskUserFieldUpdate,
		DeleteContext: resourceZendeskUserFieldDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "The URL for this user field.",
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
				Description: "The title of the user field.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"key": {
				Description: "A unique key that identifies this custom field. This is used for updating the field and referencing in placeholders. The key must consist of only letters, numbers, and underscores. It can't be only numbers and can't be reused if deleted.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Describes the purpose of the user field to users.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"position": {
				Description: "The relative position of the user field on a ticket. Note that for accounts with ticket forms, positions are controlled by the different forms.",
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
			// cannot find in user_fields doc
			// "required": {
			// 	Description: "If true, agents must enter a value in the field to change the ticket status to solved.",
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// },
			// "collapsed_for_agents": {
			// 	Description: "If true, the field is shown to agents by default. If false, the field is hidden alongside infrequently used fields. Classic interface only.",
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// },
			"regexp_for_validation": {
				Description: `For "regexp" fields only. The validation pattern for a field value to be deemed valid.`,
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				// Regular expression field only
				//TODO: validation
			},
			// "title_in_portal": {
			// 	Description: "The title of the user field for end users in Help Center.",
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Computed:    true,
			// },
			// "visible_in_portal": {
			// 	Description: "Whether this field is visible to end users in Help Center.",
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// },
			// "editable_in_portal": {
			// 	Description: "Whether this field is editable by end users in Help Center.",
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// },
			// "required_in_portal": {
			// 	Description: "If true, end users must enter a value in the field to create the request.",
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// },
			"tag": {
				Description: `For "checkbox" fields only. A tag added to tickets when the checkbox field is selected.`,
				Type:        schema.TypeString,
				Optional:    true,
			},
			// "system_field_options": {
			// 	Description: `Presented for a system user field of type "tickettype", "priority" or "status".`,
			// 	Type:        schema.TypeSet,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"name": {
			// 				Description: "System field option name.",
			// 				Type:        schema.TypeString,
			// 				Optional:    true,
			// 			},
			// 			"value": {
			// 				Description: "System field option value.",
			// 				Type:        schema.TypeString,
			// 				Optional:    true,
			// 			},
			// 		},
			// 	},
			// 	Computed: true,
			// },
			// https://developer.zendesk.com/api-reference/ticketing/tickets/user_fields/#updating-drop-down-field-options
			"custom_field_option": {
				Description: `Required and presented for a custom user field of type "dropdown". At the start set custom_field_option.id as -1, then after execution replace it with the generated id of the custom_field_option.
				Order is maintained, reorder the custom_field_option to apply the order change in dropdown in the UI`,
				Type: schema.TypeList,
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
							Required:    true,
						},
					},
				},
				Optional: true,
				//TODO: empty is invalid form
			},
			// "priority" and "status" fields only
			// "sub_type_id": {
			// 	Description: `For system user fields of type "priority" and "status". Defaults to 0. A "priority" sub type of 1 removes the "Low" and "Urgent" options. A "status" sub type of 1 adds the "On-Hold" option.`,
			// 	Type:        schema.TypeInt,
			// 	Optional:    true,
			// 	//TODO: validation
			// },
			// NOTE: Maybe this is not necessary because it's only for system field
			// "removable": {
			// 	Description: "If false, this field is a system field that must be present on all tickets.",
			// 	Type:        schema.TypeBool,
			// 	Computed:    true,
			// },
			// "agent_description": {
			// 	Description: "A description of the user field that only agents can see.",
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// },
		},
	}
}

// marshalUserField encodes the provided user field into the provided resource data
func marshalUserField(field UserField, d identifiableGetterSetter) error {
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

// unmarshalUserField parses the provided ResourceData and returns a user field
func unmarshalUserField(d identifiableGetterSetter) (UserField, error) {
	tf := UserField{}

	if v := d.Id(); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return tf, fmt.Errorf("could not parse user field id %s: %v", v, err)
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
		options := v.([]interface{})
		customFieldOptions := make([]CustomFieldOption, 0)
		for _, o := range options {
			option, ok := o.(map[string]interface{})
			if !ok {
				return tf, fmt.Errorf("could not parse custom options for field %v", tf)
			}

			optionId := option["id"]
			var idPointer *int
			if optionId != nil {
				v, ok := optionId.(int)
				if ok {
					if v == -1 {
						idPointer = nil
					} else {
						idPointer = &v
					}
				} else {
					return tf, fmt.Errorf("optionId could not be set pointer %s", optionId)
				}
			}

			customFieldOptions = append(customFieldOptions, CustomFieldOption{
				Name:  option["name"].(string),
				Value: option["value"].(string),
				ID:    idPointer,
			})
		}

		tf.CustomFieldOptions = customFieldOptions
		debugLog(tf.CustomFieldOptions, "customFieldOption")
	}

	return tf, nil
}

func resourceZendeskUserFieldCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return createUserField(ctx, d, zd)
}

func createUserField(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	tf, err := unmarshalUserField(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API request
	tf, err = CreateUserField(ctx, zd, tf)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", tf.ID))

	err = marshalUserField(tf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func CreateUserField(ctx context.Context, z *newClient.Client, userField UserField) (UserField, error) {
	var data, result struct {
		UserField UserField `json:"user_field"`
	}
	data.UserField = userField

	debugLog(data, "createUserField")

	body, err := z.Post(ctx, "/user_fields.json", data)
	if err != nil {
		return UserField{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return UserField{}, err
	}
	return result.UserField, nil
}

func resourceZendeskUserFieldRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return readUserField(ctx, d, zd)
}

func readUserField(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	field, err := GetUserField(ctx, zd, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalUserField(field, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZendeskUserFieldUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if hasChange := d.HasChange("type"); hasChange {
		if hasChange {
			return diag.FromErr(
				fmt.Errorf("field is write-once. The 'type' cannot be changed after resource creation"),
			)
		}
	}
	zd := meta.(*newClient.Client)
	return updateUserField(ctx, d, zd)
}

func updateUserField(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	tf, err := unmarshalUserField(d)
	if err != nil {
		return diag.FromErr(err)
	}

	debugLog(tf, "unmarshalled data")

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API request
	tf, err = UpdateUserField(ctx, zd, id, tf)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalUserField(tf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZendeskUserFieldDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return deleteUserField(ctx, d, zd)
}

func deleteUserField(ctx context.Context, d identifiable, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	err = DeleteUserField(ctx, zd, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// GetUserField gets a specified ticket field
// ref: https://developer.zendesk.com/rest_api/docs/support/user_fields#show-ticket-field
func GetUserField(ctx context.Context, z *newClient.Client, userID int64) (UserField, error) {
	var result struct {
		UserField UserField `json:"user_field"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/user_fields/%d.json", userID))

	if err != nil {
		return UserField{}, err
	}

	err = json.Unmarshal(body, &result)
	debugLog(result, "## GET ##")
	if err != nil {
		return UserField{}, err
	}

	return result.UserField, err
}

// UpdateUserField updates a field with the specified ticket field
// ref: https://developer.zendesk.com/rest_api/docs/support/user_fields#update-ticket-field
func UpdateUserField(ctx context.Context, z *newClient.Client, ticketID int64, field UserField) (UserField, error) {
	var result, data struct {
		UserField UserField `json:"user_field"`
	}

	data.UserField = field

	debugLog(data, "updateUserField")

	body, err := z.Put(ctx, fmt.Sprintf("/user_fields/%d.json", ticketID), data)

	if err != nil {
		return UserField{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return UserField{}, err
	}

	return result.UserField, err
}

// DeleteUserField deletes the specified ticket field
// ref: https://developer.zendesk.com/rest_api/docs/support/user_fields#Delete-ticket-field
func DeleteUserField(ctx context.Context, z *newClient.Client, ticketID int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/user_fields/%d.json", ticketID))

	if err != nil {
		return err
	}

	return nil
}
