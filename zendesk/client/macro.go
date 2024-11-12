package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nukosuke/terraform-provider-zendesk/zendesk/models"

	"github.com/nukosuke/go-zendesk/zendesk"
)

// MacroAPI an interface containing all macro related methods
type MacroAPI interface {
	GetMacros(ctx context.Context, options *zendesk.MacroListOptions) ([]models.Macro, zendesk.Page, error)
	CreateMacro(ctx context.Context, macro models.Macro) (models.Macro, error)
	DeleteMacro(ctx context.Context, id int64) error
	UpdateMacro(ctx context.Context, id int64, form models.Macro) (models.Macro, error)
	GetMacro(ctx context.Context, id int64) (models.Macro, error)
	UpdateMacroPosition(ctx context.Context, id int64, macro models.MacroPosition) error
}

// GetMacros fetches macros
// ref: https://developer.zendesk.com/rest_api/docs/support/macro#list-macros
func (z *Client) GetMacros(ctx context.Context, options *zendesk.MacroListOptions) ([]models.Macro, zendesk.Page, error) {
	var data struct {
		Macros []models.Macro `json:"macros"`
		zendesk.Page
	}

	tmp := options
	if tmp == nil {
		tmp = &zendesk.MacroListOptions{}
	}

	u, err := addOptions("/macros.json", tmp)
	if err != nil {
		return nil, zendesk.Page{}, err
	}

	body, err := z.Get(ctx, u)
	if err != nil {
		return []models.Macro{}, zendesk.Page{}, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return []models.Macro{}, zendesk.Page{}, err
	}
	return data.Macros, data.Page, nil
}

// CreateMacro creates new macro
// ref: https://developer.zendesk.com/rest_api/docs/support/macro#create-macros
func (z *Client) CreateMacro(ctx context.Context, macro models.Macro) (models.Macro, error) {
	var data, result struct {
		Macro models.Macro `json:"macro"`
	}
	data.Macro = macro

	body, err := z.Post(ctx, "/macros.json", data)
	if err != nil {
		return models.Macro{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return models.Macro{}, err
	}
	return result.Macro, nil
}

// GetMacro returns the specified macro
// ref: https://developer.zendesk.com/rest_api/docs/support/macro#show-macro
func (z *Client) GetMacro(ctx context.Context, id int64) (models.Macro, error) {
	var result struct {
		Macro models.Macro `json:"macro"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/macros/%d.json", id))
	if err != nil {
		return models.Macro{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return models.Macro{}, err
	}
	return result.Macro, nil
}

// UpdateMacro updates the specified macro and returns the updated form
// ref: https://developer.zendesk.com/rest_api/docs/support/macro#update-macros
func (z *Client) UpdateMacro(ctx context.Context, id int64, form models.Macro) (models.Macro, error) {
	var data, result struct {
		Macro models.Macro `json:"macro"`
	}

	z.UpdateMacroPosition(ctx, id, models.MacroPosition{
		ID:       id,
		Position: form.Position,
	})

	data.Macro = form
	fmt.Printf("\nUpdating Macro position to %d\n", data.Macro.Position)
	body, err := z.Put(ctx, fmt.Sprintf("/macros/%d.json", id), data)
	if err != nil {
		return models.Macro{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return models.Macro{}, err
	}

	return result.Macro, nil
}

func (z *Client) UpdateMacroPosition(ctx context.Context, id int64, macro models.MacroPosition) error {
	var data, result struct {
		Macros []models.MacroPosition `json:"macros"`
	}

	data.Macros = append(data.Macros, macro)

	body, err := z.Put(ctx, "/macros/update_many", data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	fmt.Sprintf("\nUpdated position to %d for macro %d", macro.Position, macro.ID)

	return nil
}

// DeleteMacro deletes the specified macro
// ref: https://developer.zendesk.com/rest_api/docs/support/macro#delete-macro
func (z *Client) DeleteMacro(ctx context.Context, id int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/macros/%d.json", id))
	if err != nil {
		return err
	}

	return nil
}
