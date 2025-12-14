package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// Locale represents a Zendesk locale
type Locale struct {
	ID          int64  `json:"id"`
	URL         string `json:"url"`
	Locale      string `json:"locale"`
	Name        string `json:"name"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// LocaleListResponse represents the response from listing locales
type LocaleListResponse struct {
	Locales []Locale `json:"locales"`
}

// LocaleAPI interface for locale operations
type LocaleAPI interface {
	GetLocales(ctx context.Context) ([]Locale, error)
	GetLocale(ctx context.Context, id int64) (Locale, error)
	GetAgentLocales(ctx context.Context) ([]Locale, error)
	GetCurrentLocale(ctx context.Context) (Locale, error)
	GetPublicLocales(ctx context.Context) ([]Locale, error)
	DetectBestLocale(ctx context.Context, acceptLanguage string) (Locale, error)
}

// GetLocales fetches all locales
// ref: https://developer.zendesk.com/api-reference/ticketing/account-configuration/locales/#list-locales
func (z *Client) GetLocales(ctx context.Context) ([]Locale, error) {
	var result LocaleListResponse

	body, err := z.Get(ctx, "/locales.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Locales, nil
}

// GetLocale returns a specific locale
// ref: https://developer.zendesk.com/api-reference/ticketing/account-configuration/locales/#show-locale
func (z *Client) GetLocale(ctx context.Context, id int64) (Locale, error) {
	var result struct {
		Locale Locale `json:"locale"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/locales/%d.json", id))
	if err != nil {
		return Locale{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Locale{}, err
	}

	return result.Locale, nil
}

// GetAgentLocales returns locales available to agents
// ref: https://developer.zendesk.com/api-reference/ticketing/account-configuration/locales/#list-agent-locales
func (z *Client) GetAgentLocales(ctx context.Context) ([]Locale, error) {
	var result LocaleListResponse

	body, err := z.Get(ctx, "/locales/agent.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Locales, nil
}

// GetCurrentLocale returns the current locale
// ref: https://developer.zendesk.com/api-reference/ticketing/account-configuration/locales/#show-current-locale
func (z *Client) GetCurrentLocale(ctx context.Context) (Locale, error) {
	var result struct {
		Locale Locale `json:"locale"`
	}

	body, err := z.Get(ctx, "/locales/current.json")
	if err != nil {
		return Locale{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Locale{}, err
	}

	return result.Locale, nil
}

// GetPublicLocales returns public locales
// ref: https://developer.zendesk.com/api-reference/ticketing/account-configuration/locales/#list-public-locales
func (z *Client) GetPublicLocales(ctx context.Context) ([]Locale, error) {
	var result LocaleListResponse

	body, err := z.Get(ctx, "/locales/public.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Locales, nil
}

// DetectBestLocale detects the best locale based on Accept-Language header
// ref: https://developer.zendesk.com/api-reference/ticketing/account-configuration/locales/#detect-best-locale
func (z *Client) DetectBestLocale(ctx context.Context, acceptLanguage string) (Locale, error) {
	var result struct {
		Locale Locale `json:"locale"`
	}

	url := fmt.Sprintf("/locales/detect_best_locale.json")
	if acceptLanguage != "" {
		url += "?accept_language=" + acceptLanguage
	}

	body, err := z.Get(ctx, url)
	if err != nil {
		return Locale{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Locale{}, err
	}

	return result.Locale, nil
}

