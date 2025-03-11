package zendesk

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
	"github.com/nukosuke/terraform-provider-zendesk/zendesk/models"
)

// https://developer.zendesk.com/api-reference/ticketing/business-rules/views/
func resourceZendeskView() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a user field resource.",
		CreateContext: resourceZendeskViewsCreate,
		ReadContext:   resourceZendeskViewsRead,
		UpdateContext: resourceZendeskViewsUpdate,
		DeleteContext: resourceZendeskViewsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "The URL for this user field.",
				Type:        schema.TypeString,
				Computed:    true,
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
			// Both the "all" and "any" parameter are optional, but at least one of them must be supplied
			"all": viewConditionSchema("Logical AND. All the conditions must be met."),
			"any": viewConditionSchema("Logical OR. Any condition can be met."),
			"group_title": {
				Description: "Sort or group the tickets by a column in the View columns table",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     false,
			},
			"sort_title": {
				Description: "Sort or group the tickets by a column in the View columns table",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     false,
			},
			"group_by": {
				Description: "Sort or group the tickets by a column in the View columns table",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     false,
			},
			"group_order": {
				Description: "asc or desc",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     false,
			},
			"sort_by": {
				Description: "Sort or group the tickets by a column in the View columns table",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     false,
			},
			"sort_order": {
				Description: "asc or desc",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     false,
			},
			"columns": {
				Description: "all the columns",
				Optional:    true,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

// marshalViews encodes the provided user field into the provided resource data
func marshalViews(field models.View, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":         field.URL,
		"title":       field.Title,
		"description": field.Description,
		"position":    field.Position,
		"active":      field.Active,
		"group_by":    field.Execution.GroupBy,
		"sort_by":     field.Execution.SortBy,
		"group_order": field.Execution.GroupOrder,
		"sort_order":  field.Execution.SortOrder,
	}

	if field.Restriction == nil {
		fields["restrictions"] = nil
	} else {
		var restrictions []int
		for _, col := range field.Restriction.IDs {
			restrictions = append(restrictions, col)
		}
		fields["restrictions"] = restrictions
	}

	var columns []string

	for _, col := range field.Execution.Columns {
		var _, isFloat = col.ID.(float64)
		if isFloat {
			columns = append(columns, strconv.FormatFloat(col.ID.(float64), 'f', -1, 64))
		} else {
			columns = append(columns, col.ID.(string))
		}
	}

	fields["columns"] = columns

	var alls []map[string]interface{}
	for _, v := range field.Conditions.All {
		m := map[string]interface{}{
			"field":    v.Field,
			"operator": v.Operator,
			"value":    v.Value,
		}
		alls = append(alls, m)
	}
	fields["all"] = alls

	var anys []map[string]interface{}
	for _, v := range field.Conditions.Any {
		m := map[string]interface{}{
			"field":    v.Field,
			"operator": v.Operator,
			"value":    v.Value,
		}
		anys = append(anys, m)
	}
	fields["any"] = anys

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

// unmarshalViews parses the provided ResourceData and returns a user field
func unmarshalViews(d identifiableGetterSetter) (models.View, error) {
	tf := models.View{}

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
	if v, ok := d.GetOk("group_by"); ok {
		tf.Execution.GroupBy = v.(string)
	}
	if v, ok := d.GetOk("sort_by"); ok {
		tf.Execution.SortBy = v.(string)
	}
	if v, ok := d.GetOk("group_order"); ok {
		tf.Execution.GroupOrder = v.(string)
	}
	if v, ok := d.GetOk("restrictions"); ok {
		var restrictions []int
		for _, ids := range v.(*schema.Set).List() {
			restrictions = append(restrictions, ids.(int))
		}
		tf.Restriction = &models.Restriction{}
		tf.Restriction.IDs = restrictions
		tf.Restriction.Type = "Group"
	} else {
		tf.Restriction = nil
	}
	if v, ok := d.GetOk("sort_order"); ok {
		tf.Execution.SortOrder = v.(string)
	}
	if v, ok := d.GetOk("columns"); ok {
		columns := v.([]interface{})
		c := []models.Column{}
		for _, col := range columns {
			var _, isFloat = col.(float64)
			if isFloat {
				c = append(c, models.Column{
					ID: col.(float64),
				})
			} else {
				c = append(c, models.Column{
					ID: col.(string),
				})
			}
		}
		tf.Execution.Columns = c
	}
	if v, ok := d.GetOk("all"); ok {
		allConditions := v.(*schema.Set).List()
		conditions := []models.ViewCondition{}
		for _, c := range allConditions {
			condition, ok := c.(map[string]interface{})
			if !ok {
				return tf, fmt.Errorf("could not parse 'all' conditions for view %v", tf)
			}
			conditions = append(conditions, models.ViewCondition{
				Field:    condition["field"].(string),
				Operator: condition["operator"].(string),
				Value:    condition["value"].(string),
			})
		}
		tf.Conditions.All = conditions
	}
	if v, ok := d.GetOk("any"); ok {
		anyConditions := v.(*schema.Set).List()
		conditions := []models.ViewCondition{}
		for _, c := range anyConditions {
			condition, ok := c.(map[string]interface{})
			if !ok {
				return tf, fmt.Errorf("could not parse 'any' conditions for view %v", tf)
			}
			conditions = append(conditions, models.ViewCondition{
				Field:    condition["field"].(string),
				Operator: condition["operator"].(string),
				Value:    condition["value"].(string),
			})
		}
		tf.Conditions.Any = conditions
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

	return tf, nil
}

func resourceZendeskViewsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return createViews(ctx, d, zd)
}

func createViews(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	tf, err := unmarshalViews(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API request
	tf, err = zd.CreateView(ctx, tf)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", tf.ID))

	err = marshalViews(tf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZendeskViewsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return readViews(ctx, d, zd)
}

func readViews(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	field, err := zd.GetView(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalViews(field, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZendeskViewsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return updateViews(ctx, d, zd)
}

func updateViews(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	tf, err := unmarshalViews(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API request
	tf, err = zd.UpdateView(ctx, id, tf)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalViews(tf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZendeskViewsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*newClient.Client)
	return deleteViews(ctx, d, zd)
}

func deleteViews(ctx context.Context, d identifiable, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.DeleteView(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func viewConditionSchema(desc string) *schema.Schema {
	return &schema.Schema{
		Description: desc,
		Type:        schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"field": {
					Description: "The name of a ticket field.",
					Type:        schema.TypeString,
					Required:    true,
				},
				"operator": {
					Description: "A comparison operator.",
					Type:        schema.TypeString,
					Required:    true,
				},
				"value": {
					Description: "The value of a ticket field.",
					Type:        schema.TypeString,
					Required:    true,
				},
			},
		},
		Optional: true,
	}
}
