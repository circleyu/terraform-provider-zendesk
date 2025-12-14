package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
)

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/
func resourceZendeskCustomRoles() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a custom role resource.",
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return createCustomRole(ctx, d, zd)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return readCustomRole(ctx, d, zd)
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return updateCustomRole(ctx, d, zd)
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return deleteCustomRole(ctx, d, zd)
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "The API url of this custom role.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the custom role.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the custom role.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"role_type": {
				Description: "The role type ID.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"configuration": {
				Description: "Configuration settings for the custom role.",
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        schema.TypeString,
			},
			"created_at": {
				Description: "The time the custom role was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "The time the custom role was last updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func marshalCustomRole(role newClient.CustomRole, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":         role.URL,
		"name":        role.Name,
		"description": role.Description,
		"role_type":   role.RoleType,
		"created_at":  role.CreatedAt,
		"updated_at":  role.UpdatedAt,
	}

	if role.Configuration != nil {
		fields["configuration"] = role.Configuration
	}

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalCustomRole(d identifiableGetterSetter) (newClient.CustomRole, error) {
	role := newClient.CustomRole{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return role, fmt.Errorf("could not parse custom role id %s: %v", v, err)
		}
		role.ID = id
	}

	if v, ok := d.GetOk("url"); ok {
		role.URL = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		role.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		role.Description = v.(string)
	}

	if v, ok := d.GetOk("role_type"); ok {
		role.RoleType = int64(v.(int))
	}

	if v, ok := d.GetOk("configuration"); ok {
		configMap := v.(map[string]interface{})
		role.Configuration = make(map[string]interface{})
		for k, val := range configMap {
			role.Configuration[k] = val
		}
	}

	return role, nil
}

func createCustomRole(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	role, err := unmarshalCustomRole(d)
	if err != nil {
		return diag.FromErr(err)
	}

	role, err = zd.CreateCustomRole(ctx, role)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", role.ID))

	err = marshalCustomRole(role, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readCustomRole(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	role, err := zd.GetCustomRole(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalCustomRole(role, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateCustomRole(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	role, err := unmarshalCustomRole(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	role, err = zd.UpdateCustomRole(ctx, id, role)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalCustomRole(role, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func deleteCustomRole(ctx context.Context, d identifiable, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.DeleteCustomRole(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

