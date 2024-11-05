package plans

import "time"

type Plan struct {
	ID             int64
	Name           string
	Description    string
	CreatedBy      string
	LastModifiedBy string
	IsPublished    bool
	Public         bool
	CreatedAt      time.Time
	Modified       time.Time
}

type CreatePlan struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name" validate:"required"`
	Description    string    `json:"description"`
	CreatedBy      string    `json:"created_by" validate:"required"`
	LastModifiedBy string    `json:"last_modified_by" validate:"required"`
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
	LastModifiedBy string  `json:"last_modified_by" validate:"required"`
	IsPublished    *bool   `json:"is_published,omitempty"`
	Public         *bool   `json:"public,omitempty"`
}

type SharePlanForUsers struct {
	PlanID    int64    `json:"plan_id" validate:"required"`
	UsersIDs  []string `json:"users_ids" validate:"required"`
	CreatedBy string   `json:"created_by" validate:"required"`
}

type DBPlan struct {
	ID             int64     `db:"id"`
	Name           string    `db:"name"`
	Description    string    `db:"description"`
	CreatedBy      string    `db:"created_by"`
	LastModifiedBy string    `db:"last_modified_by"`
	IsPublished    bool      `db:"is_published"`
	Public         bool      `db:"public"`
	CreatedAt      time.Time `db:"created_at"`
	Modified       time.Time `db:"modified"`
}

type DBSharePlanForUser struct {
	PlanID    int64     `db:"plan_id"`
	UserID    string    `db:"user_id"`
	CreatedBy string    `db:"created_by"`
	CreatedAt time.Time `db:"created_at"`
}
