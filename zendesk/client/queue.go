package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// Queue represents a Zendesk queue
type Queue struct {
	ID          int64                  `json:"id,omitempty"`
	URL         string                 `json:"url,omitempty"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Definition  map[string]interface{} `json:"definition,omitempty"`
	CreatedAt   string                 `json:"created_at,omitempty"`
	UpdatedAt   string                 `json:"updated_at,omitempty"`
}

// QueueListResponse represents the response from listing queues
type QueueListResponse struct {
	Queues []Queue `json:"queues"`
}

// QueueAPI interface for queue operations
type QueueAPI interface {
	GetQueues(ctx context.Context) ([]Queue, error)
	GetQueue(ctx context.Context, id int64) (Queue, error)
	CreateQueue(ctx context.Context, queue Queue) (Queue, error)
	UpdateQueue(ctx context.Context, id int64, queue Queue) (Queue, error)
	DeleteQueue(ctx context.Context, id int64) error
}

// GetQueues fetches all queues
// ref: https://developer.zendesk.com/api-reference/ticketing/queues/#list-queues
func (z *Client) GetQueues(ctx context.Context) ([]Queue, error) {
	var result QueueListResponse

	body, err := z.Get(ctx, "/queues.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Queues, nil
}

// GetQueue returns a specific queue
// ref: https://developer.zendesk.com/api-reference/ticketing/queues/#show-queue
func (z *Client) GetQueue(ctx context.Context, id int64) (Queue, error) {
	var result struct {
		Queue Queue `json:"queue"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/queues/%d.json", id))
	if err != nil {
		return Queue{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Queue{}, err
	}

	return result.Queue, nil
}

// CreateQueue creates a new queue
// ref: https://developer.zendesk.com/api-reference/ticketing/queues/#create-queue
func (z *Client) CreateQueue(ctx context.Context, queue Queue) (Queue, error) {
	var data, result struct {
		Queue Queue `json:"queue"`
	}
	data.Queue = queue

	body, err := z.Post(ctx, "/queues.json", data)
	if err != nil {
		return Queue{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Queue{}, err
	}

	return result.Queue, nil
}

// UpdateQueue updates a queue
// ref: https://developer.zendesk.com/api-reference/ticketing/queues/#update-queue
func (z *Client) UpdateQueue(ctx context.Context, id int64, queue Queue) (Queue, error) {
	var data, result struct {
		Queue Queue `json:"queue"`
	}
	data.Queue = queue

	body, err := z.Put(ctx, fmt.Sprintf("/queues/%d.json", id), data)
	if err != nil {
		return Queue{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Queue{}, err
	}

	return result.Queue, nil
}

// DeleteQueue deletes a queue
// ref: https://developer.zendesk.com/api-reference/ticketing/queues/#delete-queue
func (z *Client) DeleteQueue(ctx context.Context, id int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/queues/%d.json", id))
	if err != nil {
		return err
	}

	return nil
}

