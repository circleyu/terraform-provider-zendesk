package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// SatisfactionRating represents a Zendesk satisfaction rating
type SatisfactionRating struct {
	ID        int64  `json:"id"`
	URL       string `json:"url"`
	AssigneeID int64 `json:"assignee_id,omitempty"`
	GroupID   int64  `json:"group_id,omitempty"`
	RequesterID int64 `json:"requester_id,omitempty"`
	TicketID  int64  `json:"ticket_id,omitempty"`
	Score     string `json:"score"` // "offered", "unoffered", "good", "bad"
	Reason    string `json:"reason,omitempty"`
	ReasonID  int64  `json:"reason_id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// SatisfactionRatingListResponse represents the response from listing satisfaction ratings
type SatisfactionRatingListResponse struct {
	SatisfactionRatings []SatisfactionRating `json:"satisfaction_ratings"`
	NextPage            string                `json:"next_page,omitempty"`
	PreviousPage        string                `json:"previous_page,omitempty"`
	Count               int64                 `json:"count,omitempty"`
}

// SatisfactionRatingAPI interface for satisfaction rating operations
type SatisfactionRatingAPI interface {
	GetSatisfactionRatings(ctx context.Context) ([]SatisfactionRating, error)
	GetSatisfactionRating(ctx context.Context, id int64) (SatisfactionRating, error)
	GetSatisfactionRatingCount(ctx context.Context) (int64, error)
}

// GetSatisfactionRatings fetches all satisfaction ratings
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings/#list-satisfaction-ratings
func (z *Client) GetSatisfactionRatings(ctx context.Context) ([]SatisfactionRating, error) {
	var result SatisfactionRatingListResponse

	body, err := z.Get(ctx, "/satisfaction_ratings.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.SatisfactionRatings, nil
}

// GetSatisfactionRating returns a specific satisfaction rating
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings/#show-satisfaction-rating
func (z *Client) GetSatisfactionRating(ctx context.Context, id int64) (SatisfactionRating, error) {
	var result struct {
		SatisfactionRating SatisfactionRating `json:"satisfaction_rating"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/satisfaction_ratings/%d.json", id))
	if err != nil {
		return SatisfactionRating{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return SatisfactionRating{}, err
	}

	return result.SatisfactionRating, nil
}

// GetSatisfactionRatingCount returns the count of satisfaction ratings
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/satisfaction_ratings/#count-satisfaction-ratings
func (z *Client) GetSatisfactionRatingCount(ctx context.Context) (int64, error) {
	var result struct {
		Count int64 `json:"count"`
	}

	body, err := z.Get(ctx, "/satisfaction_ratings/count.json")
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, err
	}

	return result.Count, nil
}

