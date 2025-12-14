package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	newClient "github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
)

// https://developer.zendesk.com/api-reference/ticketing/users/users/
func resourceZendeskUsers() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a user resource.",
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return createUser(ctx, d, zd)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return readUser(ctx, d, zd)
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return updateUser(ctx, d, zd)
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*newClient.Client)
			return deleteUser(ctx, d, zd)
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "The API url of this user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"email": {
				Description: "The email address of the user.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"role": {
				Description: "The role of the user. Allowed values: end-user, agent, admin.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "end-user",
			},
			"active": {
				Description: "Whether the user is active.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"verified": {
				Description: "Whether the user is verified.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"phone": {
				Description: "The phone number of the user.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"organization_id": {
				Description: "The ID of the organization the user belongs to.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"tags": {
				Description: "Tags for the user.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Description: "The time the user was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "The time the user was last updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func marshalUser(user newClient.User, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":            user.URL,
		"name":           user.Name,
		"email":          user.Email,
		"role":           user.Role,
		"active":         user.Active,
		"verified":       user.Verified,
		"phone":          user.Phone,
		"organization_id": int(user.OrganizationID),
		"created_at":     user.CreatedAt,
		"updated_at":     user.UpdatedAt,
	}

	if user.Tags != nil {
		fields["tags"] = user.Tags
	}

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalUser(d identifiableGetterSetter) (newClient.User, error) {
	user := newClient.User{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return user, fmt.Errorf("could not parse user id %s: %v", v, err)
		}
		user.ID = id
	}

	if v, ok := d.GetOk("name"); ok {
		user.Name = v.(string)
	}

	if v, ok := d.GetOk("email"); ok {
		user.Email = v.(string)
	}

	if v, ok := d.GetOk("role"); ok {
		user.Role = v.(string)
	}

	if v, ok := d.GetOk("active"); ok {
		user.Active = v.(bool)
	}

	if v, ok := d.GetOk("phone"); ok {
		user.Phone = v.(string)
	}

	if v, ok := d.GetOk("organization_id"); ok {
		user.OrganizationID = int64(v.(int))
	}

	if v, ok := d.GetOk("tags"); ok {
		tagsSet := v.(*schema.Set)
		tags := make([]string, tagsSet.Len())
		for i, tag := range tagsSet.List() {
			tags[i] = tag.(string)
		}
		user.Tags = tags
	}

	return user, nil
}

func createUser(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	user, err := unmarshalUser(d)
	if err != nil {
		return diag.FromErr(err)
	}

	user, err = zd.CreateUser(ctx, user)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", user.ID))

	err = marshalUser(user, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readUser(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	user, err := zd.GetUser(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalUser(user, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateUser(ctx context.Context, d identifiableGetterSetter, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	user, err := unmarshalUser(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	user, err = zd.UpdateUser(ctx, id, user)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalUser(user, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func deleteUser(ctx context.Context, d identifiable, zd *newClient.Client) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.DeleteUser(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

