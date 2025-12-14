package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
)

// https://developer.zendesk.com/api-reference/ticketing/users/organization_memberships/
func resourceZendeskOrganizationMemberships() *schema.Resource {
	return &schema.Resource{
		Description: "Provides an organization membership resource.",
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return createOrganizationMembership(ctx, d, zd)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return readOrganizationMembership(ctx, d, zd)
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			// Organization memberships are immutable, so update is a no-op
			return readOrganizationMembership(ctx, d, meta.(*newClient.Client))
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return deleteOrganizationMembership(ctx, d, zd)
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "The API url of this organization membership.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"user_id": {
				Description: "The ID of the user.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"organization_id": {
				Description: "The ID of the organization.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"default": {
				Description: "Whether this is the default organization membership for the user.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"created_at": {
				Description: "The time the organization membership was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "The time the organization membership was last updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func marshalOrganizationMembership(membership newClient.OrganizationMembership, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":             membership.URL,
		"user_id":         int(membership.UserID),
		"organization_id": int(membership.OrganizationID),
		"default":         membership.Default,
		"created_at":      membership.CreatedAt,
		"updated_at":      membership.UpdatedAt,
	}

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalOrganizationMembership(d identifiableGetterSetter) (newClient.OrganizationMembership, error) {
	membership := newClient.OrganizationMembership{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return membership, fmt.Errorf("could not parse organization membership id %s: %v", v, err)
		}
		membership.ID = id
	}

	if v, ok := d.GetOk("url"); ok {
		membership.URL = v.(string)
	}

	if v, ok := d.GetOk("user_id"); ok {
		membership.UserID = int64(v.(int))
	}

	if v, ok := d.GetOk("organization_id"); ok {
		membership.OrganizationID = int64(v.(int))
	}

	if v, ok := d.GetOk("default"); ok {
		membership.Default = v.(bool)
	}

	return membership, nil
}

func createOrganizationMembership(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	membership, err := unmarshalOrganizationMembership(d)
	if err != nil {
		return diag.FromErr(err)
	}

	membership, err = zd.CreateOrganizationMembership(ctx, membership)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", membership.ID))

	err = marshalOrganizationMembership(membership, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readOrganizationMembership(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	membership, err := zd.GetOrganizationMembership(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalOrganizationMembership(membership, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func deleteOrganizationMembership(ctx context.Context, d identifiable, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.DeleteOrganizationMembership(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

