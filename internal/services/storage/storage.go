package storage

import "errors"

var (
	ErrChannelExitsts  = errors.New("channel already exists")
	ErrChannelNotFound = errors.New("channel not found")

	ErrPlanExitsts  = errors.New("plan already exists")
	ErrPlanNotFound = errors.New("plan not found")

	ErrLessonExitsts  = errors.New("lesson already exists")
	ErrLessonNotFound = errors.New("lesson not found")

	ErrPageExitsts      = errors.New("page already exists")
	ErrPageNotFound     = errors.New("page not found")
	ErrQuestionNotFound = errors.New("question not found")
	ErrUnContType       = errors.New("unsupported content type")

	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrRowsIteration         = errors.New("rows iteration failed")
	ErrScanFailed            = errors.New("scan failed")
	ErrQueryFailed           = errors.New("query failed")
	ErrFailedTransaction     = errors.New("failed to begin transaction")
	ErrFailedContTransaction = errors.New("failed to continue transaction")
	ErrRollBack              = errors.New("failed to rollback transaction")
	ErrCommitTransaction     = errors.New("failed to commit transaction")

	ErrLessonAttemtNotFound = errors.New("lesson attempt not found")
	ErrPageAttemtsNotFound  = errors.New("page attempts not found")
	ErrAnswerNotFound       = errors.New("page answer not found")

	ErrPlanAlreadySharedWithUser = errors.New("plan already shared with user")
)
