package zendesk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
)

func dataSourceZendeskSatisfactionRatings() *schema.Resource {
	return &schema.Resource{
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return readSatisfactionRatingsDataSource(ctx, d, zd)
		},

		Schema: map[string]*schema.Schema{
			"satisfaction_ratings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"assignee_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"requester_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ticket_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"score": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reason_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Description: "List of satisfaction ratings",
			},
			"count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total count of satisfaction ratings",
			},
		},
	}
}

func readSatisfactionRatingsDataSource(ctx context.Context, d *schema.ResourceData, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	ratings, err := zd.GetSatisfactionRatings(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	count, err := zd.GetSatisfactionRatingCount(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	ratingList := make([]map[string]interface{}, len(ratings))
	for i, rating := range ratings {
		ratingList[i] = map[string]interface{}{
			"id":           int(rating.ID),
			"url":          rating.URL,
			"assignee_id":  int(rating.AssigneeID),
			"group_id":     int(rating.GroupID),
			"requester_id": int(rating.RequesterID),
			"ticket_id":    int(rating.TicketID),
			"score":        rating.Score,
			"reason":       rating.Reason,
			"reason_id":    int(rating.ReasonID),
			"created_at":   rating.CreatedAt,
			"updated_at":   rating.UpdatedAt,
		}
	}

	d.SetId("satisfaction_ratings")
	d.Set("satisfaction_ratings", ratingList)
	d.Set("count", int(count))

	return diags
}

