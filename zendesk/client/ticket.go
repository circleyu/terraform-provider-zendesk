package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// Ticket represents a Zendesk ticket
type Ticket struct {
	ID              int64                  `json:"id,omitempty"`
	URL             string                 `json:"url,omitempty"`
	ExternalID      string                 `json:"external_id,omitempty"`
	Type            string                 `json:"type,omitempty"` // "question", "incident", "problem", "task"
	Subject         string                 `json:"subject"`
	RawSubject      string                 `json:"raw_subject,omitempty"`
	Description     string                 `json:"description"`
	Priority        string                 `json:"priority,omitempty"` // "urgent", "high", "normal", "low"
	Status          string                 `json:"status,omitempty"` // "new", "open", "pending", "hold", "solved", "closed"
	Recipient       string                 `json:"recipient,omitempty"`
	RequesterID     int64                  `json:"requester_id,omitempty"`
	SubmitterID     int64                  `json:"submitter_id,omitempty"`
	AssigneeID      int64                  `json:"assignee_id,omitempty"`
	OrganizationID  int64                  `json:"organization_id,omitempty"`
	GroupID         int64                  `json:"group_id,omitempty"`
	CollaboratorIDs []int64               `json:"collaborator_ids,omitempty"`
	FollowerIDs     []int64               `json:"follower_ids,omitempty"`
	EmailCCIDs      []int64               `json:"email_cc_ids,omitempty"`
	ForumTopicID    int64                  `json:"forum_topic_id,omitempty"`
	ProblemID       int64                  `json:"problem_id,omitempty"`
	HasIncidents    bool                   `json:"has_incidents,omitempty"`
	IsPublic        bool                   `json:"is_public,omitempty"`
	DueAt           string                 `json:"due_at,omitempty"`
	Tags            []string               `json:"tags,omitempty"`
	CustomFields    []map[string]interface{} `json:"custom_fields,omitempty"`
	SatisfactionRating map[string]interface{} `json:"satisfaction_rating,omitempty"`
	SharingAgreementIDs []int64             `json:"sharing_agreement_ids,omitempty"`
	Fields          []map[string]interface{} `json:"fields,omitempty"`
	FollowupIDs     []int64               `json:"followup_ids,omitempty"`
	TicketFormID    int64                  `json:"ticket_form_id,omitempty"`
	BrandID         int64                  `json:"brand_id,omitempty"`
	CreatedAt       string                 `json:"created_at,omitempty"`
	UpdatedAt       string                 `json:"updated_at,omitempty"`
}

// TicketListResponse represents the response from listing tickets
type TicketListResponse struct {
	Tickets    []Ticket `json:"tickets"`
	NextPage   string   `json:"next_page,omitempty"`
	PreviousPage string `json:"previous_page,omitempty"`
	Count      int64    `json:"count,omitempty"`
}

// TicketAPI interface for ticket operations
type TicketAPI interface {
	GetTickets(ctx context.Context) ([]Ticket, error)
	GetTicket(ctx context.Context, id int64) (Ticket, error)
	CreateTicket(ctx context.Context, ticket Ticket) (Ticket, error)
	UpdateTicket(ctx context.Context, id int64, ticket Ticket) (Ticket, error)
	DeleteTicket(ctx context.Context, id int64) error
}

// GetTickets fetches all tickets
// ref: https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#list-tickets
func (z *Client) GetTickets(ctx context.Context) ([]Ticket, error) {
	var result TicketListResponse

	body, err := z.Get(ctx, "/tickets.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Tickets, nil
}

// GetTicket returns a specific ticket
// ref: https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#show-ticket
func (z *Client) GetTicket(ctx context.Context, id int64) (Ticket, error) {
	var result struct {
		Ticket Ticket `json:"ticket"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/tickets/%d.json", id))
	if err != nil {
		return Ticket{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Ticket{}, err
	}

	return result.Ticket, nil
}

// CreateTicket creates a new ticket
// ref: https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#create-ticket
func (z *Client) CreateTicket(ctx context.Context, ticket Ticket) (Ticket, error) {
	var data, result struct {
		Ticket Ticket `json:"ticket"`
	}
	data.Ticket = ticket

	body, err := z.Post(ctx, "/tickets.json", data)
	if err != nil {
		return Ticket{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Ticket{}, err
	}

	return result.Ticket, nil
}

// UpdateTicket updates a ticket
// ref: https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#update-ticket
func (z *Client) UpdateTicket(ctx context.Context, id int64, ticket Ticket) (Ticket, error) {
	var data, result struct {
		Ticket Ticket `json:"ticket"`
	}
	data.Ticket = ticket

	body, err := z.Put(ctx, fmt.Sprintf("/tickets/%d.json", id), data)
	if err != nil {
		return Ticket{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Ticket{}, err
	}

	return result.Ticket, nil
}

// DeleteTicket deletes a ticket
// ref: https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/#delete-ticket
func (z *Client) DeleteTicket(ctx context.Context, id int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/tickets/%d.json", id))
	if err != nil {
		return err
	}

	return nil
}

