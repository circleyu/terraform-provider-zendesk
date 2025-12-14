package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
)

// https://developer.zendesk.com/api-reference/ticketing/queues/
func resourceZendeskQueues() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a queue resource.",
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return createQueue(ctx, d, zd)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return readQueue(ctx, d, zd)
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return updateQueue(ctx, d, zd)
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return deleteQueue(ctx, d, zd)
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "The API url of this queue.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the queue.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the queue.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"definition": {
				Description: "The definition of the queue (JSON object).",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        schema.TypeString,
			},
			"created_at": {
				Description: "The time the queue was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "The time the queue was last updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func marshalQueue(queue newClient.Queue, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":         queue.URL,
		"name":        queue.Name,
		"description": queue.Description,
		"created_at":  queue.CreatedAt,
		"updated_at":  queue.UpdatedAt,
	}

	if queue.Definition != nil {
		fields["definition"] = queue.Definition
	}

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalQueue(d identifiableGetterSetter) (newClient.Queue, error) {
	queue := newClient.Queue{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return queue, fmt.Errorf("could not parse queue id %s: %v", v, err)
		}
		queue.ID = id
	}

	if v, ok := d.GetOk("url"); ok {
		queue.URL = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		queue.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		queue.Description = v.(string)
	}

	if v, ok := d.GetOk("definition"); ok {
		configMap := v.(map[string]interface{})
		queue.Definition = make(map[string]interface{})
		for k, val := range configMap {
			queue.Definition[k] = val
		}
	}

	return queue, nil
}

func createQueue(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	queue, err := unmarshalQueue(d)
	if err != nil {
		return diag.FromErr(err)
	}

	queue, err = zd.CreateQueue(ctx, queue)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", queue.ID))

	err = marshalQueue(queue, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readQueue(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	queue, err := zd.GetQueue(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalQueue(queue, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateQueue(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	queue, err := unmarshalQueue(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	queue, err = zd.UpdateQueue(ctx, id, queue)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalQueue(queue, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func deleteQueue(ctx context.Context, d identifiable, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.DeleteQueue(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

