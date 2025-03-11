package models

type TicketForm struct {
	ID                 int64            `json:"id,omitempty"`
	URL                string           `json:"url,omitempty"`
	Name               string           `json:"name"`
	RawName            string           `json:"raw_name,omitempty"`
	DisplayName        string           `json:"display_name,omitempty"`
	RawDisplayName     string           `json:"raw_display_name,omitempty"`
	Position           int64            `json:"position"`
	Active             bool             `json:"active,omitempty"`
	EndUserVisible     bool             `json:"end_user_visible,omitempty"`
	Default            bool             `json:"default,omitempty"`
	TicketFieldIDs     []int64          `json:"ticket_field_ids,omitempty"`
	InAllBrands        bool             `json:"in_all_brands,omitempty"`
	RestrictedBrandIDs []int64          `json:"restricted_brand_ids,omitempty"`
	AgentConditions    []AgentCondition `json:"agent_conditions,omitempty"`
}

type AgentCondition struct {
	ParentFieldId int64         `json:"parent_field_id"`
	Value         string        `json:"value"`
	ChildFields   []ChildFields `json:"child_fields"` // eg, matching_value, matching_value_1
}

type ChildFields struct {
	Id                 int64              `json:"id"`
	IsRequired         bool               `json:"is_required"`
	RequiredOnStatuses RequiredOnStatuses `json:"required_on_statuses"`
}

type RequiredOnStatuses struct {
	Type     string   `json:"type"`     // NO_STATUSES OR SOME_STATUSES OR ALL_STATUSES
	Statuses []string `json:"statuses"` // eg, ["new", "pending", "open"]
}
