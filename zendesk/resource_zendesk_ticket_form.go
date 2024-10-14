package zendesk

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
	"github.com/nukosuke/terraform-provider-zendesk/zendesk/models"
)

// https://developer.zendesk.com/rest_api/docs/support/ticket_forms
func resourceZendeskTicketForm() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a ticket form resource.",
		CreateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(*newClient.Client)
			return createTicketForm(ctx, data, zd)
		},
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(*newClient.Client)
			return readTicketForm(ctx, data, zd)
		},
		UpdateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(*newClient.Client)
			return updateTicketForm(ctx, data, zd)
		},
		DeleteContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(*newClient.Client)
			return deleteTicketForm(ctx, data, zd)
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "URL of the ticket form.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the form.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"display_name": {
				Description: "The name of the form that is displayed to an end user.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"position": {
				Description: "The position of this form among other forms in the account, i.e. dropdown.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"active": {
				Description: "If the form is set as active.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"end_user_visible": {
				Description: "Is the form visible to the end user.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"default": {
				Description: "Is the form the default form for this account.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"ticket_field_ids": {
				Description: "ids of all ticket fields which are in this ticket form. The products use the order of the ids to show the field values in the tickets.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
			},
			"in_all_brands": {
				Description: "Is the form available for use in all brands on this account.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"restricted_brand_ids": {
				Description: "ids of all brands that this ticket form is restricted to.",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Computed: true,
			},
			"agent_conditions": {
				Description: "Array of condition sets for agent workspaces",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parent_field_id": {
							Description: "ID of the parent field",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"value": {
							Description: "value",
							Type:        schema.TypeString,
							Required:    true,
						},
						"child_fields": {
							Description: "Child Fields",
							Type:        schema.TypeSet,
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Description: "",
										Type:        schema.TypeInt,
										Required:    true,
									},
									"is_required": {
										Description: "",
										Type:        schema.TypeBool,
										Required:    true,
									},
									"required_on_statuses": {
										Description: "",
										Type:        schema.TypeSet,
										Required:    true,
										MaxItems:    1, // Ensures only one element is allowed
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Description: "",
													Required:    true,
													Type:        schema.TypeString,
												},
												"statuses": {
													Description: "",
													Optional:    true,
													Type:        schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// unmarshalTicketField parses the provided ResourceData and returns a ticket field
func unmarshalTicketForm(d identifiableGetterSetter) (models.TicketForm, error) {
	tf := models.TicketForm{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return tf, fmt.Errorf("could not parse ticket field id %s: %v", v, err)
		}
		tf.ID = id
	}

	if v, ok := d.GetOk("url"); ok {
		tf.URL = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		tf.Name = v.(string)
		tf.RawName = v.(string)
	}

	if v, ok := d.GetOk("display_name"); ok {
		tf.DisplayName = v.(string)
		tf.RawDisplayName = v.(string)
	}

	if v, ok := d.GetOk("position"); ok {
		tf.Position = int64(v.(int))
	}

	if v, ok := d.GetOk("active"); ok {
		tf.Active = v.(bool)
	}

	if v, ok := d.GetOk("end_user_visible"); ok {
		tf.EndUserVisible = v.(bool)
	}

	if v, ok := d.GetOk("default"); ok {
		tf.Default = v.(bool)
	}

	if v, ok := d.GetOk("in_all_brands"); ok {
		tf.InAllBrands = v.(bool)
	}

	if v, ok := d.GetOk("ticket_field_ids"); ok {
		ticketFieldIDs := v.([]interface{})
		for _, ticketFieldID := range ticketFieldIDs {
			tf.TicketFieldIDs = append(tf.TicketFieldIDs, int64(ticketFieldID.(int)))
		}
	}

	if v, ok := d.GetOk("restricted_brand_ids"); ok {
		brandIDs := v.(*schema.Set).List()
		for _, id := range brandIDs {
			tf.TicketFieldIDs = append(tf.RestrictedBrandIDs, int64(id.(int)))
		}
	}

	if v, ok := d.GetOk("agent_conditions"); ok {
		agent_conditions := v.(*schema.Set).List()
		for _, agent_condition := range agent_conditions {
			AgentCondition := models.AgentCondition{}
			AgentCondition.Value = agent_condition.(map[string]interface{})["value"].(string)
			AgentCondition.ParentFieldId = int64(agent_condition.(map[string]interface{})["parent_field_id"].(int))
			ChildFields := agent_condition.(map[string]interface{})["child_fields"].(*schema.Set).List()
			// for _, child_field := range agent_condition.(map[string]interface{})["child_fields"].([]interface{}) {
			for _, child_field := range ChildFields {
				ChildField := models.ChildFields{}
				ChildField.Id = int64(child_field.(map[string]interface{})["id"].(int))
				ChildField.IsRequired = child_field.(map[string]interface{})["is_required"].(bool)

				required_on_statuses := child_field.(map[string]interface{})["required_on_statuses"].(*schema.Set).List()

				if len(required_on_statuses) > 0 {
					// Extract the first element from required_on_statuses, assuming it's a slice
					statusesInterface := required_on_statuses[0].(map[string]interface{})["statuses"].(*schema.Set).List()

					// Initialize a slice to hold the converted statuses
					statuses := make([]string, len(statusesInterface))

					// Loop through the []interface{} and convert each item to a string
					for i, v := range statusesInterface {
						statuses[i] = v.(string) // Type assert to string
					}

					ChildField.RequiredOnStatuses = models.RequiredOnStatuses{
						Type:     required_on_statuses[0].(map[string]interface{})["type"].(string),
						Statuses: statuses,
					}
				}
				AgentCondition.ChildFields = append(AgentCondition.ChildFields, ChildField)
			}
			tf.AgentConditions = append(tf.AgentConditions, AgentCondition)
		}
	}

	return tf, nil
}

// marshalTicketField encodes the provided form into the provided resource data
func marshalTicketForm(f models.TicketForm, d identifiableGetterSetter) error {

	var agentConditionsList []interface{}

	for _, agentCondition := range f.AgentConditions {
		agentConditionMap := map[string]interface{}{
			"value":           agentCondition.Value,
			"parent_field_id": agentCondition.ParentFieldId,
		}

		var childFieldsList []interface{}
		for _, childField := range agentCondition.ChildFields {

			var requiredStatuses []interface{}
			requiredStatuses = append(requiredStatuses, map[string]interface{}{
				"type":     childField.RequiredOnStatuses.Type,
				"statuses": childField.RequiredOnStatuses.Statuses,
			})

			childFieldMap := map[string]interface{}{
				"id":                   childField.Id,
				"is_required":          childField.IsRequired,
				"required_on_statuses": requiredStatuses,
			}
			childFieldsList = append(childFieldsList, childFieldMap)
		}
		agentConditionMap["child_fields"] = childFieldsList
		agentConditionsList = append(agentConditionsList, agentConditionMap)
	}

	fields := map[string]interface{}{
		"url":                  f.URL,
		"name":                 f.Name,
		"display_name":         f.DisplayName,
		"position":             f.Position,
		"active":               f.Active,
		"end_user_visible":     f.EndUserVisible,
		"default":              f.Default,
		"ticket_field_ids":     f.TicketFieldIDs,
		"in_all_brands":        f.InAllBrands,
		"restricted_brand_ids": f.RestrictedBrandIDs,
		"agent_conditions":     agentConditionsList,
	}

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

func createTicketForm(ctx context.Context, d identifiableGetterSetter, zd client.TicketFormAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	tf, err := unmarshalTicketForm(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API request
	tf, err = zd.CreateTicketForm(ctx, tf)
	if err != nil {
		return diag.FromErr(err)
	}

	// Patch from created resource
	d.SetId(fmt.Sprintf("%d", tf.ID))

	err = marshalTicketForm(tf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readTicketForm(ctx context.Context, d identifiableGetterSetter, zd client.TicketFormAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tf, err := zd.GetTicketForm(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalTicketForm(tf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateTicketForm(ctx context.Context, d identifiableGetterSetter, zd client.TicketFormAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	tf, err := unmarshalTicketForm(d)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonData, err := json.Marshal(tf)
	fmt.Println("Update Processed payload: JSON")
	fmt.Println(string(jsonData))

	tf, err = zd.UpdateTicketForm(ctx, tf.ID, tf)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalTicketForm(tf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func deleteTicketForm(ctx context.Context, d identifiable, zd client.TicketFormAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.DeleteTicketForm(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
