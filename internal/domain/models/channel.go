package models

import "time"

type Channel struct {
	ID             int64
	Name           string
	Description    string
	CreatedBy      int64
	LastModifiedBy int64
	CreatedAt      time.Time
	Modified       time.Time
}

type ChannelWithPlans struct {
	ID             int64
	Name           string
	Description    string
	CreatedBy      int64
	LastModifiedBy int64
	CreatedAt      time.Time
	Modified       time.Time
	Plans          []PlanInChannel
}

type PlanInChannel struct {
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

type CreateChannel struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name" validate:"required"`
	Description    string    `json:"description"`
	CreatedBy      int64     `json:"created_by" validate:"required"`
	LastModifiedBy int64     `json:"last_modified_by" validate:"required"`
	CreatedAt      time.Time `json:"created_at"`
	Modified       time.Time `json:"modified"`
}

type UpdateChannelRequest struct {
	ID             int64   `json:"id" validate:"required"`
	Name           *string `json:"name,omitempty"`
	Description    *string `json:"description,omitempty"`
	LastModifiedBy int64   `json:"last_modified_by" validate:"required"`
}
