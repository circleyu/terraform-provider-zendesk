package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
)

// https://developer.zendesk.com/api-reference/ticketing/users/group_memberships/
func resourceZendeskGroupMemberships() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a group membership resource.",
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return createGroupMembership(ctx, d, zd)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return readGroupMembership(ctx, d, zd)
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			// Group memberships are immutable, so update is a no-op
			return readGroupMembership(ctx, d, meta.(*newClient.Client))
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return deleteGroupMembership(ctx, d, zd)
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "The API url of this group membership.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"user_id": {
				Description: "The ID of the user.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"group_id": {
				Description: "The ID of the group.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"default": {
				Description: "Whether this is the default group membership for the user.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"created_at": {
				Description: "The time the group membership was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "The time the group membership was last updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func marshalGroupMembership(membership newClient.GroupMembership, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":        membership.URL,
		"user_id":    int(membership.UserID),
		"group_id":   int(membership.GroupID),
		"default":    membership.Default,
		"created_at": membership.CreatedAt,
		"updated_at": membership.UpdatedAt,
	}

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalGroupMembership(d identifiableGetterSetter) (newClient.GroupMembership, error) {
	membership := newClient.GroupMembership{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return membership, fmt.Errorf("could not parse group membership id %s: %v", v, err)
		}
		membership.ID = id
	}

	if v, ok := d.GetOk("url"); ok {
		membership.URL = v.(string)
	}

	if v, ok := d.GetOk("user_id"); ok {
		membership.UserID = int64(v.(int))
	}

	if v, ok := d.GetOk("group_id"); ok {
		membership.GroupID = int64(v.(int))
	}

	if v, ok := d.GetOk("default"); ok {
		membership.Default = v.(bool)
	}

	return membership, nil
}

func createGroupMembership(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	membership, err := unmarshalGroupMembership(d)
	if err != nil {
		return diag.FromErr(err)
	}

	membership, err = zd.CreateGroupMembership(ctx, membership)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", membership.ID))

	err = marshalGroupMembership(membership, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readGroupMembership(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	membership, err := zd.GetGroupMembership(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalGroupMembership(membership, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func deleteGroupMembership(ctx context.Context, d identifiable, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.DeleteGroupMembership(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

