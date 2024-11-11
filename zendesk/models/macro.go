package models

import "time"

type Macro struct {
	Actions     []MacroAction `json:"actions"`
	Active      bool          `json:"active"`
	CreatedAt   time.Time     `json:"created_at,omitempty"`
	Description interface{}   `json:"description"`
	ID          int64         `json:"id,omitempty"`
	Position    int           `json:"position,omitempty"`
	Restriction interface{}   `json:"restriction"`
	Title       string        `json:"title"`
	UpdatedAt   time.Time     `json:"updated_at,omitempty"`
	URL         string        `json:"url,omitempty"`
}

type MacroAction struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}
