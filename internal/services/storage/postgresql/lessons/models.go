package lessons

import "time"

type Lesson struct {
	ID             int64
	Name           string
	Description    string
	CreatedBy      string
	LastModifiedBy string
	CreatedAt      time.Time
	Modified       time.Time
}

type CreateLesson struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name" validate:"required"`
	Description    string    `json:"description,omitempty"`
	CreatedBy      string    `json:"created_by" validate:"required"`
	LastModifiedBy string    `json:"last_modified_by" validate:"required"`
	CreatedAt      time.Time `json:"created_at"`
	Modified       time.Time `json:"modified"`
	PlanID         int64     `json:"plan_id" validate:"required"`
}

type UpdateLessonRequest struct {
	PlanID         int64  `json:"plan_id" validate:"required"`
	LessonID       int64  `json:"lesson_id" validate:"required"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	LastModifiedBy string `json:"last_modified_by" validate:"required"`
}

type DBLesson struct {
	ID             int64     `db:"id"`
	Name           string    `db:"name"`
	Description    string    `db:"description"`
	CreatedBy      string    `db:"created_by"`
	LastModifiedBy string    `db:"last_modified_by"`
	CreatedAt      time.Time `db:"created_at"`
	Modified       time.Time `db:"modified"`
}

type GetLesson struct {
	LessonID int64 `json:"lesson_id" validate:"required"`
	PlanID   int64 `json:"plan_id" validate:"required"`
}

type GetLessons struct {
	PlanID int64 `json:"plan_id" validate:"required"`
	Limit  int64 `json:"limit,omitempty" validate:"min=1"`
	Offset int64 `json:"offset,omitempty" validate:"min=0"`
}

func (p *GetLessons) SetDefaults() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}

type DeleteLesson struct {
	LessonID int64 `json:"lesson_id" validate:"required"`
	PlanID   int64 `json:"plan_id" validate:"required"`
}
