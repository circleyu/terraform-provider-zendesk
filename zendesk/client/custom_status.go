package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// CustomStatus represents a Zendesk custom status
type CustomStatus struct {
	ID             int64  `json:"id,omitempty"`
	URL            string `json:"url,omitempty"`
	StatusCategory string `json:"status_category"` // "new", "open", "pending", "hold", "solved"
	AgentLabel     string `json:"agent_label"`
	EndUserLabel   string `json:"end_user_label,omitempty"`
	Description    string `json:"description,omitempty"`
	Default        bool   `json:"default,omitempty"`
	Active         bool   `json:"active,omitempty"`
	EndUserHidden  bool   `json:"end_user_hidden,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

// CustomStatusListResponse represents the response from listing custom statuses
type CustomStatusListResponse struct {
	CustomStatuses []CustomStatus `json:"custom_statuses"`
}

// CustomStatusAPI interface for custom status operations
type CustomStatusAPI interface {
	GetCustomStatuses(ctx context.Context, statusCategory, active, defaultStatus *string) ([]CustomStatus, error)
	GetCustomStatus(ctx context.Context, id int64) (CustomStatus, error)
	CreateCustomStatus(ctx context.Context, status CustomStatus) (CustomStatus, error)
	UpdateCustomStatus(ctx context.Context, id int64, status CustomStatus) (CustomStatus, error)
}

// GetCustomStatuses fetches all custom statuses
// ref: https://developer.zendesk.com/api-reference/ticketing/tickets/ticket-statuses/#list-custom-ticket-statuses
func (z *Client) GetCustomStatuses(ctx context.Context, statusCategory, active, defaultStatus *string) ([]CustomStatus, error) {
	var result CustomStatusListResponse

	url := "/custom_statuses.json"
	queryParams := make(map[string]string)
	if statusCategory != nil {
		queryParams["status_categories"] = *statusCategory
	}
	if active != nil {
		queryParams["active"] = *active
	}
	if defaultStatus != nil {
		queryParams["default"] = *defaultStatus
	}

	if len(queryParams) > 0 {
		url += "?"
		first := true
		for k, v := range queryParams {
			if !first {
				url += "&"
			}
			url += fmt.Sprintf("%s=%s", k, v)
			first = false
		}
	}

	body, err := z.Get(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.CustomStatuses, nil
}

// GetCustomStatus returns a specific custom status
// ref: https://developer.zendesk.com/api-reference/ticketing/tickets/ticket-statuses/#show-custom-ticket-status
func (z *Client) GetCustomStatus(ctx context.Context, id int64) (CustomStatus, error) {
	var result struct {
		CustomStatus CustomStatus `json:"custom_status"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/custom_statuses/%d.json", id))
	if err != nil {
		return CustomStatus{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return CustomStatus{}, err
	}

	return result.CustomStatus, nil
}

// CreateCustomStatus creates a new custom status
// ref: https://developer.zendesk.com/api-reference/ticketing/tickets/ticket-statuses/#create-custom-ticket-status
func (z *Client) CreateCustomStatus(ctx context.Context, status CustomStatus) (CustomStatus, error) {
	var data, result struct {
		CustomStatus CustomStatus `json:"custom_status"`
	}
	data.CustomStatus = status

	body, err := z.Post(ctx, "/custom_statuses.json", data)
	if err != nil {
		return CustomStatus{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return CustomStatus{}, err
	}

	return result.CustomStatus, nil
}

// UpdateCustomStatus updates a custom status
// ref: https://developer.zendesk.com/api-reference/ticketing/tickets/ticket-statuses/#update-custom-ticket-status
func (z *Client) UpdateCustomStatus(ctx context.Context, id int64, status CustomStatus) (CustomStatus, error) {
	var data, result struct {
		CustomStatus CustomStatus `json:"custom_status"`
	}
	data.CustomStatus = status

	body, err := z.Put(ctx, fmt.Sprintf("/custom_statuses/%d.json", id), data)
	if err != nil {
		return CustomStatus{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return CustomStatus{}, err
	}

	return result.CustomStatus, nil
}

