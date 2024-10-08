package channels

import (
	"database/sql"
	"time"
)

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

type DBChannel struct {
	ID             int64     `db:"id"`
	Name           string    `db:"name"`
	Description    string    `db:"description"`
	CreatedBy      int64     `db:"created_by"`
	LastModifiedBy int64     `db:"last_modified_by"`
	CreatedAt      time.Time `db:"created_at"`
	Modified       time.Time `db:"modified"`
}

type DBChannelWithPlans struct {
	ID             int64              `db:"id"`
	Name           string             `db:"name"`
	Description    string             `db:"description"`
	CreatedBy      int64              `db:"created_by"`
	LastModifiedBy int64              `db:"last_modified_by"`
	CreatedAt      time.Time          `db:"created_at"`
	Modified       time.Time          `db:"modified"`
	Plans          []DBPlanInChannels `db:"plans"`
}

// We can use DBPlan, but for flexibility and control returning fields init new struct
type DBPlanInChannels struct {
	ID             sql.NullInt64  `db:"id"`
	Name           sql.NullString `db:"name"`
	Description    sql.NullString `db:"description"`
	CreatedBy      sql.NullInt64  `db:"created_by"`
	LastModifiedBy sql.NullInt64  `db:"last_modified_by"`
	IsPublished    sql.NullBool   `db:"is_published"`
	Public         sql.NullBool   `db:"public"`
	CreatedAt      sql.NullTime   `db:"created_at"`
	Modified       sql.NullTime   `db:"modified"`
}
