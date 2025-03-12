package redis

type SavePageAttempt struct {
	LessonAttemptID int64
	PageAttemptID   int64
	UserAnswer      string
	IsCorrect       bool
}

type GetPagesAttempts struct {
	LessonAttemptID int64
}

type PageAttempt struct {
	QPAttemptID     int64
	LessonAttemptID int64
}

type UpdatePageAttempt struct {
	LessonAttemptID int64
	PageID          int64
	QPAttemptID     int64
	UserAnswer      string
}
