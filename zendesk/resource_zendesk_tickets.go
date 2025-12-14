package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
)

// https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/
func resourceZendeskTickets() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a ticket resource.",
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return createTicket(ctx, d, zd)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return readTicket(ctx, d, zd)
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return updateTicket(ctx, d, zd)
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return deleteTicket(ctx, d, zd)
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "The API url of this ticket.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"subject": {
				Description: "The subject of the ticket.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the ticket.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"type": {
				Description: "The type of the ticket. Allowed values: question, incident, problem, task.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"priority": {
				Description: "The priority of the ticket. Allowed values: urgent, high, normal, low.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"status": {
				Description: "The status of the ticket. Allowed values: new, open, pending, hold, solved, closed.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "new",
			},
			"requester_id": {
				Description: "The ID of the requester.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"assignee_id": {
				Description: "The ID of the assignee.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"organization_id": {
				Description: "The ID of the organization.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"group_id": {
				Description: "The ID of the group.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"tags": {
				Description: "Tags for the ticket.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Description: "The time the ticket was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "The time the ticket was last updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func marshalTicket(ticket newClient.Ticket, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":            ticket.URL,
		"subject":        ticket.Subject,
		"description":    ticket.Description,
		"type":           ticket.Type,
		"priority":      ticket.Priority,
		"status":         ticket.Status,
		"requester_id":   int(ticket.RequesterID),
		"assignee_id":    int(ticket.AssigneeID),
		"organization_id": int(ticket.OrganizationID),
		"group_id":       int(ticket.GroupID),
		"created_at":     ticket.CreatedAt,
		"updated_at":     ticket.UpdatedAt,
	}

	if ticket.Tags != nil {
		fields["tags"] = ticket.Tags
	}

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalTicket(d identifiableGetterSetter) (newClient.Ticket, error) {
	ticket := newClient.Ticket{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return ticket, fmt.Errorf("could not parse ticket id %s: %v", v, err)
		}
		ticket.ID = id
	}

	if v, ok := d.GetOk("subject"); ok {
		ticket.Subject = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		ticket.Description = v.(string)
	}

	if v, ok := d.GetOk("type"); ok {
		ticket.Type = v.(string)
	}

	if v, ok := d.GetOk("priority"); ok {
		ticket.Priority = v.(string)
	}

	if v, ok := d.GetOk("status"); ok {
		ticket.Status = v.(string)
	}

	if v, ok := d.GetOk("requester_id"); ok {
		ticket.RequesterID = int64(v.(int))
	}

	if v, ok := d.GetOk("assignee_id"); ok {
		ticket.AssigneeID = int64(v.(int))
	}

	if v, ok := d.GetOk("organization_id"); ok {
		ticket.OrganizationID = int64(v.(int))
	}

	if v, ok := d.GetOk("group_id"); ok {
		ticket.GroupID = int64(v.(int))
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsSet := v.(*schema.Set)
		tags := make([]string, tagsSet.Len())
		for i, tag := range tagsSet.List() {
			tags[i] = tag.(string)
		}
		ticket.Tags = tags
	}

	return ticket, nil
}

func createTicket(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	ticket, err := unmarshalTicket(d)
	if err != nil {
		return diag.FromErr(err)
	}

	ticket, err = zd.CreateTicket(ctx, ticket)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", ticket.ID))

	err = marshalTicket(ticket, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readTicket(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	ticket, err := zd.GetTicket(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalTicket(ticket, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateTicket(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	ticket, err := unmarshalTicket(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	ticket, err = zd.UpdateTicket(ctx, id, ticket)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalTicket(ticket, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func deleteTicket(ctx context.Context, d identifiable, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.DeleteTicket(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

