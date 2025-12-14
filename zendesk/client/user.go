package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// User represents a Zendesk user
type User struct {
	ID              int64                  `json:"id,omitempty"`
	URL             string                 `json:"url,omitempty"`
	Name            string                 `json:"name"`
	Email           string                 `json:"email,omitempty"`
	CreatedAt       string                 `json:"created_at,omitempty"`
	UpdatedAt       string                 `json:"updated_at,omitempty"`
	Active          bool                   `json:"active,omitempty"`
	Verified        bool                   `json:"verified,omitempty"`
	Shared          bool                   `json:"shared,omitempty"`
	LocaleID        int64                  `json:"locale_id,omitempty"`
	Locale          string                 `json:"locale,omitempty"`
	TimeZone        string                 `json:"time_zone,omitempty"`
	LastLoginAt     string                 `json:"last_login_at,omitempty"`
	Phone           string                 `json:"phone,omitempty"`
	Signature       string                 `json:"signature,omitempty"`
	Details         string                 `json:"details,omitempty"`
	Notes           string                 `json:"notes,omitempty"`
	OrganizationID  int64                  `json:"organization_id,omitempty"`
	Role            string                 `json:"role,omitempty"` // "end-user", "agent", "admin"
	CustomRoleID    int64                  `json:"custom_role_id,omitempty"`
	Moderator       bool                   `json:"moderator,omitempty"`
	TicketRestriction string               `json:"ticket_restriction,omitempty"`
	OnlyPrivateComments bool               `json:"only_private_comments,omitempty"`
	Tags            []string               `json:"tags,omitempty"`
	ExternalID      string                 `json:"external_id,omitempty"`
	Alias           string                 `json:"alias,omitempty"`
	UserFields      map[string]interface{} `json:"user_fields,omitempty"`
}

// UserListResponse represents the response from listing users
type UserListResponse struct {
	Users     []User `json:"users"`
	NextPage  string `json:"next_page,omitempty"`
	PreviousPage string `json:"previous_page,omitempty"`
	Count     int64  `json:"count,omitempty"`
}

// UserAPI interface for user operations
type UserAPI interface {
	GetUsers(ctx context.Context) ([]User, error)
	GetUser(ctx context.Context, id int64) (User, error)
	CreateUser(ctx context.Context, user User) (User, error)
	UpdateUser(ctx context.Context, id int64, user User) (User, error)
	DeleteUser(ctx context.Context, id int64) error
}

// GetUsers fetches all users
// ref: https://developer.zendesk.com/api-reference/ticketing/users/users/#list-users
func (z *Client) GetUsers(ctx context.Context) ([]User, error) {
	var result UserListResponse

	body, err := z.Get(ctx, "/users.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Users, nil
}

// GetUser returns a specific user
// ref: https://developer.zendesk.com/api-reference/ticketing/users/users/#show-user
func (z *Client) GetUser(ctx context.Context, id int64) (User, error) {
	var result struct {
		User User `json:"user"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/users/%d.json", id))
	if err != nil {
		return User{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return User{}, err
	}

	return result.User, nil
}

// CreateUser creates a new user
// ref: https://developer.zendesk.com/api-reference/ticketing/users/users/#create-user
func (z *Client) CreateUser(ctx context.Context, user User) (User, error) {
	var data, result struct {
		User User `json:"user"`
	}
	data.User = user

	body, err := z.Post(ctx, "/users.json", data)
	if err != nil {
		return User{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return User{}, err
	}

	return result.User, nil
}

// UpdateUser updates a user
// ref: https://developer.zendesk.com/api-reference/ticketing/users/users/#update-user
func (z *Client) UpdateUser(ctx context.Context, id int64, user User) (User, error) {
	var data, result struct {
		User User `json:"user"`
	}
	data.User = user

	body, err := z.Put(ctx, fmt.Sprintf("/users/%d.json", id), data)
	if err != nil {
		return User{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return User{}, err
	}

	return result.User, nil
}

// DeleteUser deletes a user
// ref: https://developer.zendesk.com/api-reference/ticketing/users/users/#delete-user
func (z *Client) DeleteUser(ctx context.Context, id int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/users/%d.json", id))
	if err != nil {
		return err
	}

	return nil
}

