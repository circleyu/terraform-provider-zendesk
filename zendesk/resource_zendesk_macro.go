package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
	"github.com/nukosuke/terraform-provider-zendesk/zendesk/models"
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
							Description: "The new value of the field. Can be a single string value or a jsonencode'ed list",
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
				Description: "IMPORTANT! in order for this to take affect an update on the resource is necessary, since only that triggers a call to update position",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
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
func marshalMacros(field models.Macro, d identifiableGetterSetter) error {
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

		// If the macro	action value is a string, leave it be
		// If it's a list, marshal it to a string
		var stringVal string
		switch action.Value.(type) {
		case []interface{}:
			tmp, err := json.Marshal(action.Value)
			if err != nil {
				return fmt.Errorf("error decoding macro action value: %s", err)
			}
			stringVal = string(tmp)
		case string:
			stringVal = action.Value.(string)
		}

		m := map[string]interface{}{
			"field": action.Field,
			"value": stringVal,
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
func unmarshalMacros(d identifiableGetterSetter) (models.Macro, error) {
	tf := models.Macro{}

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
		actions := []models.MacroAction{}
		for _, a := range macroActions {
			action, ok := a.(map[string]interface{})
			if !ok {
				return tf, fmt.Errorf("could not parse actions for macro %v", tf)
			}

			// If the action value is a list, unmarshal it
			var actionValue interface{}
			if strings.HasPrefix(action["value"].(string), "[") {
				err := json.Unmarshal([]byte(action["value"].(string)), &actionValue)
				if err != nil {
					return tf, fmt.Errorf("error unmarshalling macro action value: %s", err)
				}
			} else {
				actionValue = action["value"]
			}

			actions = append(actions, models.MacroAction{
				Field: action["field"].(string),
				Value: actionValue,
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
