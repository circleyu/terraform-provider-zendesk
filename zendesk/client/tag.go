package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// Tag represents a Zendesk tag
type Tag struct {
	Name string `json:"name"`
	Count int64  `json:"count,omitempty"`
}

// TagListResponse represents the response from listing tags
type TagListResponse struct {
	Tags  []Tag `json:"tags"`
	Count int64 `json:"count,omitempty"`
}

// TagAPI interface for tag operations
type TagAPI interface {
	GetTags(ctx context.Context) ([]Tag, error)
	GetTagCount(ctx context.Context) (int64, error)
	AutocompleteTags(ctx context.Context, name string) ([]Tag, error)
}

// GetTags fetches all tags
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#list-tags
func (z *Client) GetTags(ctx context.Context) ([]Tag, error) {
	var result TagListResponse

	body, err := z.Get(ctx, "/tags.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Tags, nil
}

// GetTagCount returns the count of tags
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#count-tags
func (z *Client) GetTagCount(ctx context.Context) (int64, error) {
	var result struct {
		Count int64 `json:"count"`
	}

	body, err := z.Get(ctx, "/tags/count.json")
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, err
	}

	return result.Count, nil
}

// AutocompleteTags searches for tags by name
// ref: https://developer.zendesk.com/api-reference/ticketing/ticket-management/tags/#autocomplete-tags
func (z *Client) AutocompleteTags(ctx context.Context, name string) ([]Tag, error) {
	var result struct {
		Tags []Tag `json:"tags"`
	}

	url := fmt.Sprintf("/autocomplete/tags.json?name=%s", name)
	body, err := z.Get(ctx, url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Tags, nil
}

