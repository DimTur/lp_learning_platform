package questions

import (
	"time"
)

type QuestionPage struct {
	ID             int64
	LessonID       int64
	CreatedBy      int64
	LastModifiedBy int64
	CreatedAt      time.Time
	Modified       time.Time
	ContentType    string

	QuestionType string

	Question string
	OptionA  string
	OptionB  string
	OptionC  string
	OptionD  string
	OptionE  string
	Answer   string
}

type CreateQuestionPage struct {
	LessonID       int64  `json:"lesson_id" validate:"required"`
	CreatedBy      int64  `json:"created_by" validate:"required"`
	LastModifiedBy int64  `json:"last_modified_by" validate:"required"`
	ContentType    string `json:"content_type" validate:"required"`

	QuestionType string `json:"question_type" validate:"required"`

	Question string `json:"question" validate:"required"`
	OptionA  string `json:"option_a" validate:"required"`
	OptionB  string `json:"option_b" validate:"required"`
	OptionC  string `json:"option_c,omitempty"`
	OptionD  string `json:"option_d,omitempty"`
	OptionE  string `json:"option_e,omitempty"`
	Answer   string `json:"answer" validate:"required"`
}

type UpdateQuestionPage struct {
	ID             int64 `json:"id" validate:"required"`
	LastModifiedBy int64 `json:"last_modified_by" validate:"required"`

	Question *string `json:"question,omitempty"`
	OptionA  *string `json:"option_a,omitempty"`
	OptionB  *string `json:"option_b,omitempty"`
	OptionC  *string `json:"option_c,omitempty"`
	OptionD  *string `json:"option_d,omitempty"`
	OptionE  *string `json:"option_e,omitempty"`
	Answer   *string `json:"answer,omitempty"`
}

type DBQuestionPage struct {
	ID             int64     `db:"id"`
	LessonID       int64     `db:"lesson_id"`
	CreatedBy      int64     `db:"created_by"`
	LastModifiedBy int64     `db:"last_modified_by"`
	CreatedAt      time.Time `db:"created_at"`
	Modified       time.Time `db:"modified"`
	ContentType    string    `db:"content_type"`

	QuestionType string `db:"question_type"`

	Question string `db:"question"`
	OptionA  string `db:"option_a"`
	OptionB  string `db:"option_b"`
	OptionC  string `db:"option_c"`
	OptionD  string `db:"option_d"`
	OptionE  string `db:"option_e"`
	Answer   string `db:"answer"`
}
