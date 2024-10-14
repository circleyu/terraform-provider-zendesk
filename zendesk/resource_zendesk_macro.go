package zendesk

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/nukosuke/go-zendesk/zendesk"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
)

// https://developer.zendesk.com/api-reference/ticketing/business-rules/macros/
func resourceZendeskMacro() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a user field resource.",
		CreateContext: resourceZendeskMacrosCreate,
		ReadContext:   resourceZendeskMacrosRead,
		UpdateContext: resourceZendeskMacrosUpdate,
		DeleteContext: resourceZendeskMacrosDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "The URL for this user field.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"action": {
				Description: "What the macro will do.",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Description: "The name of a ticket field to modify.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"value": {
							Description: "The new value of the field.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
				Required: true,
			},
			"title": {
				Description: "The title of the user field.",
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
			"restrictions": {
				Description: "allowed group ids",
				Optional:    true,
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

// marshalMacros encodes the provided user field into the provided resource data
func marshalMacros(field client.Macro, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":         field.URL,
		"title":       field.Title,
		"description": field.Description,
		"position":    field.Position,
		"active":      field.Active,
	}

	if field.Restriction == nil {
		fields["restrictions"] = nil
	} else {
		var restrictions []int
		mapi := field.Restriction.(map[string]interface{})
		fmt.Println("marshalling to terraform")
		fmt.Println(mapi)
		ids := mapi["ids"]
		if ids == nil {
			fields["restrictions"] = nil
		} else {
			for _, col := range ids.([]interface{}) {
				restrictions = append(restrictions, int(col.(float64)))
			}
			fields["restrictions"] = restrictions

		}
	}

	var actions []map[string]interface{}
	for _, action := range field.Actions {

		// If the field value is a string, leave it be
		// If it's a list, marshal it to a string
		// var stringVal string
		// switch action.Value.(type) {
		// case []interface{}:
		// 	tmp, err := json.Marshal(action.Value)
		// 	if err != nil {
		// 		return fmt.Errorf("error decoding field action value: %s", err)
		// 	}
		// 	stringVal = string(tmp)
		// case string:
		// 	stringVal = action.Value.(string)
		// }

		m := map[string]interface{}{
			"field": action.Field,
			"value": action.Value,
		}
		actions = append(actions, m)
	}
	fields["action"] = actions

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

// unmarshalMacros parses the provided ResourceData and returns a user field
func unmarshalMacros(d identifiableGetterSetter) (client.Macro, error) {
	tf := client.Macro{}

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

	if v, ok := d.GetOk("title"); ok {
		tf.Title = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		tf.Description = v.(string)
	}

	if v, ok := d.GetOk("position"); ok {
		tf.Position = v.(int)
	}

	if v, ok := d.GetOk("active"); ok {
		tf.Active = v.(bool)
	}

	if v, ok := d.GetOk("restrictions"); ok {
		var restrictions []int
		for _, ids := range v.(*schema.Set).List() {
			restrictions = append(restrictions, ids.(int))
		}
		macroRestriction := &MacroRestriction{}
		macroRestriction.IDs = restrictions
		macroRestriction.Type = "Group"

		tf.Restriction = macroRestriction
	} else {
		tf.Restriction = nil
	}

	if v, ok := d.GetOk("action"); ok {
		macroActions := v.(*schema.Set).List()
		actions := []client.MacroAction{}
		for _, a := range macroActions {
			action, ok := a.(map[string]interface{})
			if !ok {
				return tf, fmt.Errorf("could not parse actions for macro %v", tf)
			}

			actions = append(actions, client.MacroAction{
				Field: action["field"].(string),
				Value: action["value"].(string),
			})
		}
		tf.Actions = actions
	}

	return tf, nil
}

func resourceZendeskMacrosCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return createMacros(ctx, d, zd)
}

func createMacros(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	tf, err := unmarshalMacros(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API request
	tf, err = zd.CreateMacro(ctx, tf)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", tf.ID))

	err = marshalMacros(tf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZendeskMacrosRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return readMacros(ctx, d, zd)
}

func readMacros(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	field, err := zd.GetMacro(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalMacros(field, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZendeskMacrosUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return updateMacros(ctx, d, zd)
}

func updateMacros(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	tf, err := unmarshalMacros(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API request
	tf, err = zd.UpdateMacro(ctx, id, tf)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalMacros(tf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZendeskMacrosDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return deleteMacros(ctx, d, zd)
}

func deleteMacros(ctx context.Context, d identifiable, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.DeleteMacro(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

type MacroRestriction struct {
	IDs  []int  `json:"ids"`
	Type string `json:"type"`
}
