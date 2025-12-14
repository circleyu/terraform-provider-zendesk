package client

import (
	"context"
	"encoding/json"
)

// OAuthClient represents a Zendesk OAuth client
type OAuthClient struct {
	ID          int64  `json:"id,omitempty"`
	URL         string `json:"url,omitempty"`
	Name        string `json:"name"`
	Identifier  string `json:"identifier,omitempty"`
	Secret      string `json:"secret,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// OAuthClientListResponse represents the response from listing OAuth clients
type OAuthClientListResponse struct {
	Clients []OAuthClient `json:"clients"`
}

// OAuthClientAPI interface for OAuth client operations
type OAuthClientAPI interface {
	GetOAuthClients(ctx context.Context) ([]OAuthClient, error)
}

// GetOAuthClients fetches all OAuth clients
// ref: https://developer.zendesk.com/api-reference/ticketing/account-configuration/oauth_clients/#list-oauth-clients
func (z *Client) GetOAuthClients(ctx context.Context) ([]OAuthClient, error) {
	var result OAuthClientListResponse

	body, err := z.Get(ctx, "/oauth/clients.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Clients, nil
}

