package attempts

type CreateAttempt interface {
	GetCommonFields()
	GetContentTypeSpecificFields() []interface{}
	GetInsertQuery() string
}

type CreateLessonAttempt struct {
	LessonID  int64 `json:"lesson_id" validate:"required"`
	PlanId    int64 `json:"plan_id" validate:"required"`
	ChannelID int64 `json:"channel_id" validate:"required"`
	UserID    int64 `json:"user_id" validate:"required"`
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
	PageID        int64 `json:"page_id" validate:"required"`
	PageAttemptID int64 `json:"page_attempt_id" validate:"required"`
}

type QuestionPage struct {
	ContentType    string `json:"content_type" validate:"required"`
	QuestionType   string `json:"question_type" validate:"required"`
	QuestionPageID int64  `json:"question_questionpage_id" validate:"required"`
}

type DBQuestionPage struct {
	ContentType    string `db:"content_type"`
	QuestionType   string `db:"question_type"`
	QuestionPageID int64  `db:"question_questionpage_id"`
}
