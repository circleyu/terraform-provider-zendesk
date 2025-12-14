package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// CustomRole represents a Zendesk custom role
type CustomRole struct {
	ID          int64    `json:"id,omitempty"`
	URL         string   `json:"url,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	RoleType    int64    `json:"role_type,omitempty"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}

// CustomRoleListResponse represents the response from listing custom roles
type CustomRoleListResponse struct {
	CustomRoles []CustomRole `json:"custom_roles"`
}

// CustomRoleAPI interface for custom role operations
type CustomRoleAPI interface {
	GetCustomRoles(ctx context.Context) ([]CustomRole, error)
	GetCustomRole(ctx context.Context, id int64) (CustomRole, error)
	CreateCustomRole(ctx context.Context, role CustomRole) (CustomRole, error)
	UpdateCustomRole(ctx context.Context, id int64, role CustomRole) (CustomRole, error)
	DeleteCustomRole(ctx context.Context, id int64) error
}

// GetCustomRoles fetches all custom roles
// ref: https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/#list-custom-roles
func (z *Client) GetCustomRoles(ctx context.Context) ([]CustomRole, error) {
	var result CustomRoleListResponse

	body, err := z.Get(ctx, "/custom_roles.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.CustomRoles, nil
}

// GetCustomRole returns a specific custom role
// ref: https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/#show-custom-role
func (z *Client) GetCustomRole(ctx context.Context, id int64) (CustomRole, error) {
	var result struct {
		CustomRole CustomRole `json:"custom_role"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/custom_roles/%d.json", id))
	if err != nil {
		return CustomRole{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return CustomRole{}, err
	}

	return result.CustomRole, nil
}

// CreateCustomRole creates a new custom role
// ref: https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/#create-custom-role
func (z *Client) CreateCustomRole(ctx context.Context, role CustomRole) (CustomRole, error) {
	var data, result struct {
		CustomRole CustomRole `json:"custom_role"`
	}
	data.CustomRole = role

	body, err := z.Post(ctx, "/custom_roles.json", data)
	if err != nil {
		return CustomRole{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return CustomRole{}, err
	}

	return result.CustomRole, nil
}

// UpdateCustomRole updates a custom role
// ref: https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/#update-custom-role
func (z *Client) UpdateCustomRole(ctx context.Context, id int64, role CustomRole) (CustomRole, error) {
	var data, result struct {
		CustomRole CustomRole `json:"custom_role"`
	}
	data.CustomRole = role

	body, err := z.Put(ctx, fmt.Sprintf("/custom_roles/%d.json", id), data)
	if err != nil {
		return CustomRole{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return CustomRole{}, err
	}

	return result.CustomRole, nil
}

// DeleteCustomRole deletes a custom role
// ref: https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/#delete-custom-role
func (z *Client) DeleteCustomRole(ctx context.Context, id int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/custom_roles/%d.json", id))
	if err != nil {
		return err
	}

	return nil
}

