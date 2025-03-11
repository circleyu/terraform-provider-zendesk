package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nukosuke/terraform-provider-zendesk/zendesk/models"

	"github.com/nukosuke/go-zendesk/zendesk"
)

// TicketFormAPI an interface containing all ticket form related methods
type TicketFormAPI interface {
	GetTicketForms(ctx context.Context, options *zendesk.TicketFormListOptions) ([]models.TicketForm, zendesk.Page, error)
	CreateTicketForm(ctx context.Context, ticketForm models.TicketForm) (models.TicketForm, error)
	DeleteTicketForm(ctx context.Context, id int64) error
	UpdateTicketForm(ctx context.Context, id int64, form models.TicketForm) (models.TicketForm, error)
	GetTicketForm(ctx context.Context, id int64) (models.TicketForm, error)
}

// GetTicketForms fetches ticket forms
// ref: https://developer.zendesk.com/rest_api/docs/support/ticket_forms#list-ticket-forms
func (z *Client) GetTicketForms(ctx context.Context, options *zendesk.TicketFormListOptions) ([]models.TicketForm, zendesk.Page, error) {
	var data struct {
		TicketForms []models.TicketForm `json:"ticket_forms"`
		zendesk.Page
	}

	tmp := options
	if tmp == nil {
		tmp = &zendesk.TicketFormListOptions{}
	}

	u, err := addOptions("/ticket_forms.json", tmp)
	if err != nil {
		return nil, zendesk.Page{}, err
	}

	body, err := z.Get(ctx, u)
	if err != nil {
		return []models.TicketForm{}, zendesk.Page{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return []models.TicketForm{}, zendesk.Page{}, err
	}
	return data.TicketForms, data.Page, nil
}

// CreateTicketForm creates new ticket form
// ref: https://developer.zendesk.com/rest_api/docs/support/ticket_forms#create-ticket-forms
func (z *Client) CreateTicketForm(ctx context.Context, ticketForm models.TicketForm) (models.TicketForm, error) {
	var data, result struct {
		TicketForm models.TicketForm `json:"ticket_form"`
	}
	data.TicketForm = ticketForm

	body, err := z.Post(ctx, "/ticket_forms.json", data)
	if err != nil {
		return models.TicketForm{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return models.TicketForm{}, err
	}
	return result.TicketForm, nil
}

// GetTicketForm returns the specified ticket form
// ref: https://developer.zendesk.com/rest_api/docs/support/ticket_forms#show-ticket-form
func (z *Client) GetTicketForm(ctx context.Context, id int64) (models.TicketForm, error) {
	var result struct {
		TicketForm models.TicketForm `json:"ticket_form"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/ticket_forms/%d.json", id))
	if err != nil {
		return models.TicketForm{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return models.TicketForm{}, err
	}
	return result.TicketForm, nil
}

// UpdateTicketForm updates the specified ticket form and returns the updated form
// ref: https://developer.zendesk.com/rest_api/docs/support/ticket_forms#update-ticket-forms
func (z *Client) UpdateTicketForm(ctx context.Context, id int64, form models.TicketForm) (models.TicketForm, error) {
	var data, result struct {
		TicketForm models.TicketForm `json:"ticket_form"`
	}

	data.TicketForm = form
	body, err := z.Put(ctx, fmt.Sprintf("/ticket_forms/%d.json", id), data)
	if err != nil {
		return models.TicketForm{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return models.TicketForm{}, err
	}

	return result.TicketForm, nil
}

// DeleteTicketForm deletes the specified ticket form
// ref: https://developer.zendesk.com/rest_api/docs/support/ticket_forms#delete-ticket-form
func (z *Client) DeleteTicketForm(ctx context.Context, id int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/ticket_forms/%d.json", id))
	if err != nil {
		return err
	}

	return nil
}
