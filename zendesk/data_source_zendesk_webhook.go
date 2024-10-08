package zendesk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/nukosuke/go-zendesk/zendesk"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
)

func dataSourceZendeskWebhook() *schema.Resource {
	return &schema.Resource{
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(*newClient.Client)
			return readWebhookDataSource(ctx, data, zd)
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"authentication": {
				Description: "Adds authentication to the webhook's HTTP requests.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Type of authentication.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"data": {
							Description: "Authentication data.",
							Type:        schema.TypeMap,
							Computed:    true,
							Sensitive:   true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"add_position": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"description": {
				Description: "Webhook description.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"endpoint": {
				Description: "The destination URL that the webhook notifies when Zendesk events occur.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"http_method": {
				Description: `HTTP method used for the webhook's requests. To subscribe the webhook to Zendesk events, this must be "POST". Allowed values are "GET", "POST", "PUT", "PATCH", or "DELETE".`,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Webhook name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"request_format": {
				Description: `The format of the data that the webhook will send. To subscribe the webhook to Zendesk events, this must be "json". Allowed values are "json", "xml", or "form_encoded".`,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status": {
				Description: `Current status of the webhook. Allowed values are "active", or "inactive".`,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"subscriptions": {
				Description: `Event subscriptions for the webhook. To subscribe the webhook to Zendesk events, specify one or more event types. For supported event type values, see Webhook event types. To connect the webhook to a trigger or automation, specify only "conditional_ticket_events" in the array.`,
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_at": {
				Description: "When the webhook was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_by": {
				Description: `ID of the user who created the webhook. "-1" represents the Zendesk system.`,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "When the webhook was updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_by": {
				Description: `ID of the user who last updated the webhook.`,
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func readWebhookDataSource(ctx context.Context, d identifiableGetterSetter, zd client.WebhookAPI) diag.Diagnostics {
	id := d.Get("id").(string)
	d.SetId(id)
	var diags diag.Diagnostics

	wh, err := zd.GetWebhook(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalWebhookDataSource(wh, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func marshalWebhookDataSource(wh *client.Webhook, d identifiableGetterSetter) error {
	fields := map[string]any{
		"description":    wh.Description,
		"endpoint":       wh.Endpoint,
		"http_method":    wh.HTTPMethod,
		"name":           wh.Name,
		"request_format": wh.RequestFormat,
		"status":         wh.Status,
		"subscriptions":  wh.Subscriptions,
		"created_at":     wh.CreatedAt.String(),
		"created_by":     wh.CreatedBy,
		"updated_at":     wh.UpdatedAt.String(),
		"updated_by":     wh.UpdatedBy,
	}

	if wh.Authentication != nil {
		auth := map[string]any{
			"type":         wh.Authentication.Type,
			"data":         wh.Authentication.Data.(map[string]any),
			"add_position": wh.Authentication.AddPosition,
		}
		fields["authentication"] = []map[string]any{auth}
	}

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}
