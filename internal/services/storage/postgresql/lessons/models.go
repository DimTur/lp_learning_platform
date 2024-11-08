package lessons

import "time"

type Lesson struct {
	ID             int64
	Name           string
	CreatedBy      int64
	LastModifiedBy int64
	CreatedAt      time.Time
	Modified       time.Time
}

type CreateLesson struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name" validate:"required"`
	CreatedBy      int64     `json:"created_by" validate:"required"`
	LastModifiedBy int64     `json:"last_modified_by" validate:"required"`
	CreatedAt      time.Time `json:"created_at"`
	Modified       time.Time `json:"modified"`
	PlanID         int64     `json:"plan_id" validate:"required"`
}

type UpdateLessonRequest struct {
	ID             int64   `json:"id" validate:"required"`
	Name           *string `json:"name,omitempty"`
	LastModifiedBy int64   `json:"last_modified_by" validate:"required"`
}

type DBLesson struct {
	ID             int64     `db:"id"`
	Name           string    `db:"name"`
	CreatedBy      int64     `db:"created_by"`
	LastModifiedBy int64     `db:"last_modified_by"`
	CreatedAt      time.Time `db:"created_at"`
	Modified       time.Time `db:"modified"`
}
