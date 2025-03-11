package models

import "time"

type (
	// View is struct for group membership payload
	// https://developer.zendesk.com/api-reference/ticketing/business-rules/views/

	ViewCondition struct {
		Field    string `json:"field"`
		Operator string `json:"operator"`
		Value    string `json:"value"`
	}

	Column struct {
		ID    interface{} `json:"id"` // is string for normal fields & number for custom fields
		Title string      `json:"title"`
	}

	Restriction struct {
		IDs  []int  `json:"ids"`
		Type string `json:"type"`
	}

	// View has a certain structure in Get & Different structure in
	// Put/Post
	View struct {
		ID          int64        `json:"id,omitempty"`
		Active      bool         `json:"active"`
		Description string       `json:"description"`
		Position    int          `json:"position"`
		Title       string       `json:"title"`
		CreatedAt   time.Time    `json:"created_at,omitempty"`
		UpdatedAt   time.Time    `json:"updated_at,omitempty"`
		Restriction *Restriction `json:"restriction"`
		Conditions  struct {
			All []ViewCondition `json:"all"`
			Any []ViewCondition `json:"any"`
		} `json:"conditions"`
		URL       string `json:"url,omitempty"`
		Execution struct {
			Columns    []Column `json:"columns"`
			GroupBy    string   `json:"group_by,omitempty"`
			SortBy     string   `json:"sort_by,omitempty"`
			GroupOrder string   `json:"group_order,omitempty"`
			SortOrder  string   `json:"sort_order,omitempty"`
		} `json:"execution"`
	}
	ViewCreateOrUpdate struct {
		ID          int64           `json:"id,omitempty"`
		Active      bool            `json:"active"`
		Description string          `json:"description"`
		Position    int             `json:"position"`
		Title       string          `json:"title"`
		CreatedAt   time.Time       `json:"created_at,omitempty"`
		UpdatedAt   time.Time       `json:"updated_at,omitempty"`
		All         []ViewCondition `json:"all"`
		Any         []ViewCondition `json:"any"`
		URL         string          `json:"url,omitempty"`
		Restriction *Restriction    `json:"restriction"`

		Output struct {
			Columns    []interface{} `json:"columns"` // number for custom fields, string otherwise
			GroupBy    string        `json:"group_by,omitempty"`
			SortBy     string        `json:"sort_by,omitempty"`
			GroupOrder string        `json:"group_order,omitempty"`
			SortOrder  string        `json:"sort_order,omitempty"`
		} `json:"output"`
	}

	ViewPosition struct {
		ID       int64 `json:"id,omitempty"`
		Position int   `json:"position,omitempty"`
	}
)
