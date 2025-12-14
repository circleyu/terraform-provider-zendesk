package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nukosuke/terraform-provider-zendesk/zendesk/models"
)

type ViewAPI interface {
	CreateView(ctx context.Context, view models.View) (models.Macro, error)
	DeleteView(ctx context.Context, id int64) error
	UpdateView(ctx context.Context, id int64, form models.Macro) (models.Macro, error)
	GetView(ctx context.Context, id int64) (models.Macro, error)
	UpdateViewPosition(ctx context.Context, id int64, view models.ViewPosition) error
}

func mapViewToViewCreateOrUpdate(view models.View) models.ViewCreateOrUpdate {
	var viewCreateOrUpdate models.ViewCreateOrUpdate

	// Map properties from view to viewCreateOrUpdate
	viewCreateOrUpdate.ID = view.ID
	viewCreateOrUpdate.Active = view.Active
	viewCreateOrUpdate.Description = view.Description
	viewCreateOrUpdate.Position = view.Position
	viewCreateOrUpdate.Title = view.Title
	viewCreateOrUpdate.CreatedAt = view.CreatedAt
	viewCreateOrUpdate.UpdatedAt = view.UpdatedAt
	viewCreateOrUpdate.All = view.Conditions.All
	viewCreateOrUpdate.Any = view.Conditions.Any
	viewCreateOrUpdate.URL = view.URL

	// Rename "Execution" to "Output" in ViewCreateOrUpdate
	viewCreateOrUpdate.Output.GroupBy = view.Execution.GroupBy
	viewCreateOrUpdate.Output.SortBy = view.Execution.SortBy
	viewCreateOrUpdate.Output.GroupOrder = view.Execution.GroupOrder
	viewCreateOrUpdate.Output.SortOrder = view.Execution.SortOrder

	viewCreateOrUpdate.Restriction = view.Restriction

	var columns []interface{}
	for _, col := range view.Execution.Columns {
		columns = append(columns, col.ID)
	}
	viewCreateOrUpdate.Output.Columns = columns

	return viewCreateOrUpdate
}

func (z *Client) CreateView(ctx context.Context, view models.View) (models.View, error) {
	var result struct {
		View models.View `json:"view"`
	}
	var data struct {
		View models.ViewCreateOrUpdate `json:"view"`
	}
	data.View = mapViewToViewCreateOrUpdate(view)

	body, err := z.Post(ctx, "/views.json", data)

	if err != nil {
		return models.View{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return models.View{}, err
	}
	return result.View, nil
}

func (z *Client) GetView(ctx context.Context, id int64) (models.View, error) {
	var result struct {
		View models.View `json:"view"`
	}

	body, err := z.Get(ctx, fmt.Sprintf("/views/%d.json", id))
	fmt.Println("GET bar")
	fmt.Println(string(body))

	if err != nil {
		return models.View{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return models.View{}, err
	}

	return result.View, err
}

// UpdateView updates a field with the specified ticket field
// ref: https://developer.zendesk.com/rest_api/docs/support/user_fields#update-ticket-field
func (z *Client) UpdateView(ctx context.Context, id int64, view models.View) (models.View, error) {
	var result struct {
		View models.View `json:"view"`
	}
	var data struct {
		View models.ViewCreateOrUpdate `json:"view"`
	}

	data.View = mapViewToViewCreateOrUpdate(view)

	jsonData, err := json.Marshal(data)
	fmt.Println("Update Processed payload: JSON")
	fmt.Println(string(jsonData))

	z.UpdateViewPosition(ctx, id, models.ViewPosition{
		ID:       id,
		Position: view.Position,
	})

	body, err := z.Put(ctx, fmt.Sprintf("/views/%d.json", id), data)

	if err != nil {
		fmt.Println("Printing Error")
		fmt.Printf("%+v\n", err)
		return models.View{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return models.View{}, err
	}

	return result.View, err
}

func (z *Client) UpdateViewPosition(ctx context.Context, id int64, view models.ViewPosition) error {
	var data, result struct {
		Views []models.ViewPosition `json:"views"`
	}

	data.Views = append(data.Views, view)

	body, err := z.Put(ctx, "/views/update_many", data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}

	fmt.Sprintf("\nUpdated position to %d for view %d", view.Position, view.ID)

	return nil
}

// DeleteView deletes the specified ticket field
// ref: https://developer.zendesk.com/rest_api/docs/support/user_fields#Delete-ticket-field
func (z *Client) DeleteView(ctx context.Context, viewID int64) error {
	err := z.Delete(ctx, fmt.Sprintf("/views/%d.json", viewID))

	if err != nil {
		return err
	}

	return nil
}
