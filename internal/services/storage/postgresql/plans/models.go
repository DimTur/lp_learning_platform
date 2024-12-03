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

type GetPlan struct {
	PlanID    int64 `json:"plan_id" validate:"required"`
	ChannelID int64 `json:"channel_id" validate:"required"`
}

type DeletePlan struct {
	PlanID    int64 `json:"plan_id" validate:"required"`
	ChannelID int64 `json:"channel_id" validate:"required"`
}

type UpdatePlanRequest struct {
	ChannelID      int64  `json:"channel_id" validate:"required"`
	PlanID         int64  `json:"plan_id" validate:"required"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	LastModifiedBy string `json:"last_modified_by" validate:"required"`
	IsPublished    bool   `json:"is_published,omitempty"`
	Public         bool   `json:"public,omitempty"`
}

type SharePlanForUsers struct {
	ChannelID int64     `json:"channel_id" validate:"required"`
	PlanID    int64     `json:"plan_id" validate:"required"`
	UserIDs   []string  `json:"user_ids" validate:"required"`
	CreatedBy string    `json:"created_by" validate:"required"`
	CreatedAt time.Time `json:"created_at,omitempty"`
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
	ChannelID int64     `db:"channel_id"`
	PlanID    int64     `db:"plan_id"`
	UserID    string    `db:"user_id"`
	CreatedBy string    `db:"created_by"`
	CreatedAt time.Time `db:"created_at"`
}

type BatchSharePlan struct {
	PlanID    int64
	UserIDs   []string
	CreatedBy string
	CreatedAt time.Time
}

type DBCanShare struct {
	ChannelID int64 `db:"channel_id"`
	PlanID    int64 `db:"plan_id"`
}

type IsUserShareWithPlan struct {
	UserID string `json:"user_id" validate:"required"`
	PlanID int64  `json:"plan_id" validate:"required"`
}

type GetPlans struct {
	UserID    string `json:"user_id" validate:"required"`
	ChannelID int64  `json:"channel_id" validate:"required"`
	Limit     int64  `json:"limit,omitempty" validate:"min=1"`
	Offset    int64  `json:"offset,omitempty" validate:"min=0"`
}

func (p *GetPlans) SetDefaults() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}

type LearningGroup struct {
	LgID string `json:"learning_group_id" validate:"required"`
}

type DBPlansForSharing struct {
	ChannelID int64 `db:"channel_id"`
	PlanID    int64 `db:"plan_id"`
}

type PlansForSharing struct {
	ChannelID int64
	PlanID    int64
}
