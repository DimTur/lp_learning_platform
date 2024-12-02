package attempts

import "time"

type CreateAttempt interface {
	GetCommonFields()
	GetContentTypeSpecificFields() []interface{}
	GetInsertQuery() string
}

type CreateLessonAttempt struct {
	LessonID  int64     `json:"lesson_id" validate:"required"`
	PlanID    int64     `json:"plan_id" validate:"required"`
	ChannelID int64     `json:"channel_id" validate:"required"`
	UserID    string    `json:"user_id" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
}

type CreateAbstractPageAttempt struct {
	LessonAttemptID int64  `json:"lesson_attempt_id" validate:"required"`
	ContentType     string `json:"content_type" validate:"required"`
}

type CreateAbstractQuestionAttempt struct {
	QuestionType  string `json:"question_type" validate:"required"`
	PageAttemptID int64  `json:"page_attempt_id" validate:"required"`
}

type CreateQuestionPageAttempt struct {
	PageID            int64 `json:"page_id" validate:"required"`
	QuestionAttemptID int64 `json:"question_attempt_id" validate:"required"`
}

type QuestionPage struct {
	ContentType    string `json:"content_type" validate:"required"`
	QuestionType   string `json:"question_type" validate:"required"`
	QuestionPageID int64  `json:"question_questionpage_id" validate:"required"`
}

type CreateQuestionPageAttemptNew struct {
	CreateAbstractPageAttempt
	CreateAbstractQuestionAttempt
	CreateQuestionPageAttempt
	Modified time.Time `json:"modified" validate:"required"`
}

type CreateQuestionPageAttemptResp struct {
	ID     int64 `json:"id"`
	PageID int64 `json:"page_id"`
}

type UpdatePageAttempt struct {
	QPAttemptID  int64     `json:"question_page_attempt_id" validate:"required"`
	UserAnswer   string    `json:"user_answer,omitempty"`
	Modified     time.Time `json:"modified" validate:"required"`
	IsSuccessful bool      `json:"is_successful" validate:"required"`
}

type UpdateLessonAttempt struct {
	UserID          string    `json:"user_id" validate:"required"`
	LessonAttemptID int64     `json:"lesson_attempt_id" validate:"required"`
	EndTime         time.Time `json:"end_time" validate:"required"`
	IsComplete      bool      `json:"is_complete" validate:"required"`
	IsSuccessful    bool      `json:"is_successful" validate:"required"`
	PercentageScore int64     `json:"percentage_score" validate:"required"`
}

type GetQuestionPageAttempts struct {
	UserID    string `json:"user_id" validate:"required"`
	LessonID  int64  `json:"lesson_id" validate:"required"`
	PlanID    int64  `json:"plan_id" validate:"required"`
	ChannelID int64  `json:"channel_id" validate:"required"`
}

type QuestionPageAttempt struct {
	ID              int64  `json:"id" redis:"id"`
	PageID          int64  `json:"page_id" redis:"page_id"`
	LessonAttemptID int64  `json:"lesson_attempt_id" redis:"lesson_attempt_id"`
	IsCorrect       bool   `json:"is_correct" redis:"is_correct"`
	UserAnswer      string `json:"user_answer,omitempty"`
}

type QuestionPageAttemptNew struct {
	ID              int64 `json:"id"`
	PageID          int64 `json:"page_id"`
	LessonAttemptID int64 `json:"lesson_attempt_id"`
}

type DBQuestionPageAttempt struct {
	ID              int64  `db:"id"`
	PageID          int64  `db:"page_id"`
	LessonAttemptID int64  `db:"lesson_attempt_id"`
	IsCorrect       bool   `db:"is_correct"`
	UserAnswer      string `db:"user_answer"`
}

type DBAnswer struct {
	Answer string `db:"answer"`
}

type DBQuestionPage struct {
	ContentType    string `db:"content_type"`
	QuestionType   string `db:"question_type"`
	QuestionPageID int64  `db:"question_questionpage_id"`
}

type CompleteLessonRequest struct {
	UserID          string `json:"user_id" validate:"required"`
	LessonAttemptID int64  `json:"lesson_attempt_id" validate:"required"`
}

type CompleteLessonResp struct {
	ID              int64
	IsSuccessful    bool
	PercentageScore int64
}

type GetLessonAttempts struct {
	UserID   string `json:"user_id" validate:"required"`
	LessonID int64  `json:"lesson_id,omitempty"`
	Limit    int64  `json:"limit,omitempty" validate:"min=1"`
	Offset   int64  `json:"offset,omitempty" validate:"min=0"`
}

func (p *GetLessonAttempts) SetDefaults() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}

type LessonAttempt struct {
	ID              int64
	UserID          string
	LessonID        int64
	PlanID          int64
	ChannelID       int64
	StartTime       time.Time
	EndTime         time.Time
	IsComplete      bool
	IsSuccessful    bool
	PercentageScore int64
}

type DBLessonAttempt struct {
	ID              int64     `db:"id"`
	UserID          string    `db:"user_id"`
	LessonID        int64     `db:"lesson_id"`
	PlanID          int64     `db:"plan_id"`
	ChannelID       int64     `db:"channel_id"`
	StartTime       time.Time `db:"start_time"`
	EndTime         time.Time `db:"end_time"`
	IsComplete      bool      `db:"is_complete"`
	IsSuccessful    bool      `db:"is_successful"`
	PercentageScore int64     `db:"percentage_score"`
}

type GetLessonAttemptsResp struct {
	LessonAttempts []LessonAttempt
}

type PermissionForUser struct {
	UserID          string `json:"user_id" validate:"required"`
	LessonAttemptID int64  `json:"lesson_attempt_id" validate:"required"`
}
