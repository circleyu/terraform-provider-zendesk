package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// GroupMembership represents a Zendesk group membership
type GroupMembership struct {
	ID        int64  `json:"id,omitempty"`
	URL       string `json:"url,omitempty"`
	UserID    int64  `json:"user_id"`
	GroupID   int64  `json:"group_id"`
	Default   bool   `json:"default,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// GroupMembershipListResponse represents the response from listing group memberships
type GroupMembershipListResponse struct {
	GroupMemberships []GroupMembership `json:"group_memberships"`
	NextPage         string            `json:"next_page,omitempty"`
	PreviousPage     string            `json:"previous_page,omitempty"`
	Count            int64             `json:"count,omitempty"`
}

// OrganizationMembership represents a Zendesk organization membership
type OrganizationMembership struct {
	ID             int64  `json:"id,omitempty"`
	URL            string `json:"url,omitempty"`
	UserID         int64  `json:"user_id"`
	OrganizationID int64  `json:"organization_id"`
	Default        bool   `json:"default,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

// OrganizationMembershipListResponse represents the response from listing organization memberships
type OrganizationMembershipListResponse struct {
	OrganizationMemberships []OrganizationMembership `json:"organization_memberships"`
	NextPage                string                   `json:"next_page,omitempty"`
	PreviousPage            string                   `json:"previous_page,omitempty"`
	Count                   int64                   `json:"count,omitempty"`
}

// GroupMembershipAPI interface for group membership operations
type GroupMembershipAPI interface {
	GetGroupMemberships(ctx context.Context) ([]GroupMembership, error)
	GetGroupMembership(ctx context.Context, id int64) (GroupMembership, error)
	CreateGroupMembership(ctx context.Context, membership GroupMembership) (GroupMembership, error)
	DeleteGroupMembership(ctx context.Context, id int64) error
}

// OrganizationMembershipAPI interface for organization membership operations
type OrganizationMembershipAPI interface {
	GetOrganizationMemberships(ctx context.Context) ([]OrganizationMembership, error)
	GetOrganizationMembership(ctx context.Context, id int64) (OrganizationMembership, error)
	CreateOrganizationMembership(ctx context.Context, membership OrganizationMembership) (OrganizationMembership, error)
	DeleteOrganizationMembership(ctx context.Context, id int64) error
}

// GetGroupMemberships fetches all group memberships
// ref: https://developer.zendesk.com/api-reference/ticketing/users/group_memberships/#list-group-memberships
func (z *Client) GetGroupMemberships(ctx context.Context) ([]GroupMembership, error) {
	var result GroupMembershipListResponse

	body, err := z.Get(ctx, "/group_memberships.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.GroupMemberships, nil
}

// GetGroupMembership returns a specific group membership
// ref: https://developer.zendesk.com/api-reference/ticketing/users/group_memberships/#show-group-membership
func (z *Client) GetGroupMembership(ctx context.Context, id int64) (GroupMembership, error) {
	var result struct {
		GroupMembership GroupMembership `json:"group_membership"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/group_memberships/%d.json", id))
	if err != nil {
		return GroupMembership{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return GroupMembership{}, err
	}

	return result.GroupMembership, nil
}

// CreateGroupMembership creates a new group membership
// ref: https://developer.zendesk.com/api-reference/ticketing/users/group_memberships/#create-group-membership
func (z *Client) CreateGroupMembership(ctx context.Context, membership GroupMembership) (GroupMembership, error) {
	var data, result struct {
		GroupMembership GroupMembership `json:"group_membership"`
	}
	data.GroupMembership = membership

	body, err := z.Post(ctx, "/group_memberships.json", data)
	if err != nil {
		return GroupMembership{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return GroupMembership{}, err
	}

	return result.GroupMembership, nil
}

// DeleteGroupMembership deletes a group membership
// ref: https://developer.zendesk.com/api-reference/ticketing/users/group_memberships/#delete-group-membership
func (z *Client) DeleteGroupMembership(ctx context.Context, id int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/group_memberships/%d.json", id))
	if err != nil {
		return err
	}

	return nil
}

// GetOrganizationMemberships fetches all organization memberships
// ref: https://developer.zendesk.com/api-reference/ticketing/users/organization_memberships/#list-organization-memberships
func (z *Client) GetOrganizationMemberships(ctx context.Context) ([]OrganizationMembership, error) {
	var result OrganizationMembershipListResponse

	body, err := z.Get(ctx, "/organization_memberships.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.OrganizationMemberships, nil
}

// GetOrganizationMembership returns a specific organization membership
// ref: https://developer.zendesk.com/api-reference/ticketing/users/organization_memberships/#show-organization-membership
func (z *Client) GetOrganizationMembership(ctx context.Context, id int64) (OrganizationMembership, error) {
	var result struct {
		OrganizationMembership OrganizationMembership `json:"organization_membership"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/organization_memberships/%d.json", id))
	if err != nil {
		return OrganizationMembership{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return OrganizationMembership{}, err
	}

	return result.OrganizationMembership, nil
}

// CreateOrganizationMembership creates a new organization membership
// ref: https://developer.zendesk.com/api-reference/ticketing/users/organization_memberships/#create-organization-membership
func (z *Client) CreateOrganizationMembership(ctx context.Context, membership OrganizationMembership) (OrganizationMembership, error) {
	var data, result struct {
		OrganizationMembership OrganizationMembership `json:"organization_membership"`
	}
	data.OrganizationMembership = membership

	body, err := z.Post(ctx, "/organization_memberships.json", data)
	if err != nil {
		return OrganizationMembership{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return OrganizationMembership{}, err
	}

	return result.OrganizationMembership, nil
}

// DeleteOrganizationMembership deletes an organization membership
// ref: https://developer.zendesk.com/api-reference/ticketing/users/organization_memberships/#delete-organization-membership
func (z *Client) DeleteOrganizationMembership(ctx context.Context, id int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/organization_memberships/%d.json", id))
	if err != nil {
		return err
	}

	return nil
}

