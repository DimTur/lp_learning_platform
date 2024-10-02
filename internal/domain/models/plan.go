package models

import "time"

type Plan struct {
	ID             int64
	Name           string
	Description    string
	CreatedBy      int64
	LastModifiedBy int64
	IsPublished    bool
	Public         bool
	CreatedAt      time.Time
	Modified       time.Time
}

type CreatePlan struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name" validate:"required"`
	Description    string    `json:"description"`
	CreatedBy      int64     `json:"created_by" validate:"required"`
	LastModifiedBy int64     `json:"last_modified_by" validate:"required"`
	IsPublished    bool      `json:"is_published"`
	Public         bool      `json:"public"`
	CreatedAt      time.Time `json:"created_at"`
	Modified       time.Time `json:"modified"`
	ChannelID      int64     `json:"channel_id" validate:"required"`
}

type UpdatePlanRequest struct {
	ID             int64   `json:"id" validate:"required"`
	Name           *string `json:"name,omitempty"`
	Description    *string `json:"description,omitempty"`
	LastModifiedBy int64   `json:"last_modified_by" validate:"required"`
	IsPublished    *bool   `json:"is_published,omitempty"`
	Public         *bool   `json:"public,omitempty"`
}
