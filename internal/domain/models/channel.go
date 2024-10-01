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

type CreateChannel struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name" validate:"required"`
	Description    string    `json:"description"`
	CreatedBy      int64     `json:"created_by" validate:"required"`
	LastModifiedBy int64     `json:"last_modified_by" validate:"required"`
	CreatedAt      time.Time `json:"created_at"`
	Modified       time.Time `json:"modified"`
}

type DBChannel struct {
	ID             int64     `db:"id"`
	Name           string    `db:"name"`
	Description    string    `db:"description"`
	CreatedBy      int64     `db:"created_by"`
	LastModifiedBy int64     `db:"last_modified_by"`
	CreatedAt      time.Time `db:"created_at"`
	Modified       time.Time `db:"modified"`
}

type UpdateChannelRequest struct {
	ID             int64   `json:"id" validate:"required"`
	Name           *string `json:"name,omitempty"`
	Description    *string `json:"description,omitempty"`
	LastModifiedBy int64   `json:"last_modified_by" validate:"required"`
}
