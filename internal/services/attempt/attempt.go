package attempt

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/DimTur/lp_learning_platform/internal/services/redis"
	"github.com/DimTur/lp_learning_platform/internal/services/storage"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/attempts"
	"github.com/go-playground/validator/v10"
)

type AttemptSaver interface {
	CreateLessonAttempt(ctx context.Context, lAttempt *attempts.CreateLessonAttempt) (int64, error)
	CreateQuestionPageAttempts(ctx context.Context, attempt attempts.CreateQuestionPageAttemptNew) (*attempts.CreateQuestionPageAttemptResp, error)
	UpdatePageAttempt(ctx context.Context, updPAttempt *attempts.UpdatePageAttempt) error
	UpdateLessonAttempt(ctx context.Context, updLAttempt *attempts.UpdateLessonAttempt) (int64, error)
}

type AttemptProvider interface {
	GetQuestionPages(ctx context.Context, lessonID int64) ([]attempts.QuestionPage, error)
	GetLessonPagesAttempts(ctx context.Context, questionPage *attempts.GetQuestionPageAttempts) ([]attempts.QuestionPageAttempt, error)
	GetCurrentAnswerForAttempt(ctx context.Context, pageID int64) (string, error)
	CheckLessonAttempt(ctx context.Context, lessonAttempt *attempts.GetQuestionPageAttempts) (int64, error)
	GetLessonAttempts(ctx context.Context, input *attempts.GetLessonAttempts) (*attempts.GetLessonAttemptsResp, error)
	CheckPermissionForUser(ctx context.Context, userAtt *attempts.PermissionForUser) (bool, error)
}

type AttemptRedisStore interface {
	SavePageAttempt(ctx context.Context, pageAttempt *redis.SavePageAttempt) error
	GetPageAttempts(ctx context.Context, lessonAttemptID int64) ([]attempts.QuestionPageAttempt, error)
}

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrInvalidAttemptID     = errors.New("invalid attempt id")
	ErrAttemptExitsts       = errors.New("attempt already exists")
	ErrAttemptNotFound      = errors.New("attempt not found")
	ErrFailedToCreate       = errors.New("attempt creation failed")
	ErrPageAttemtsNotFound  = errors.New("page attempts not found")
	ErrFailedToSaveInRedis  = errors.New("failed to save in redis")
	ErrAnswerNotFound       = errors.New("page answer not found")
	ErrLessonAttemtNotFound = errors.New("lesson attempt not found")
	ErrPermissionsDenied    = errors.New("permissions denied")
)

type AttemptHandlers struct {
	log               *slog.Logger
	validator         *validator.Validate
	attemptSaver      AttemptSaver
	attemptProvider   AttemptProvider
	attemptRedisStore AttemptRedisStore
}

func New(
	log *slog.Logger,
	validator *validator.Validate,
	attemptSaver AttemptSaver,
	attemptProvider AttemptProvider,
	attemptRedisStore AttemptRedisStore,
) *AttemptHandlers {
	return &AttemptHandlers{
		log:               log,
		validator:         validator,
		attemptSaver:      attemptSaver,
		attemptProvider:   attemptProvider,
		attemptRedisStore: attemptRedisStore,
	}
}

// TryLesson get relevant for user questions page attempts thrue lesson attempt
// If not exist create it and return
func (ah *AttemptHandlers) TryLesson(ctx context.Context, questionPage *attempts.GetQuestionPageAttempts) ([]attempts.QuestionPageAttempt, error) {
	const op = "attempts.TryLesson"
	log := ah.log.With(
		slog.String("op", op),
		slog.String("user_id", questionPage.UserID),
		slog.Int64("lesson_id", questionPage.LessonID),
	)

	// Validation
	if err := ah.validator.Struct(questionPage); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	// Step 1: Check for an existing lesson attempt
	lessonAttemptID, err := ah.attemptProvider.CheckLessonAttempt(ctx, questionPage)
	if err != nil {
		log.Error("failed to check lesson attempt", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if lessonAttemptID != 0 {
		// Step 2: Try to get attempts from Redis or DB
		return ah.getExistingPageAttempts(ctx, lessonAttemptID, questionPage, log)
	}

	// Step 3: Create a new lesson attempt and question page attempts
	return ah.createLessonAttemptAndPages(ctx, questionPage, log)
}

func (ah *AttemptHandlers) UpdatePageAttempt(ctx context.Context, updPAttempt *redis.UpdatePageAttempt) error {
	const op = "attempts.CreateAttempt"

	log := ah.log.With(
		slog.String("op", op),
		slog.Int64("lesson_attempt_id", updPAttempt.QPAttemptID),
	)

	// Validation
	err := ah.validator.Struct(updPAttempt)
	if err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("updating attempt")

	// Get current answer
	answer, err := ah.attemptProvider.GetCurrentAnswerForAttempt(ctx, updPAttempt.PageID)
	if err != nil {
		if errors.Is(err, storage.ErrAnswerNotFound) {
			ah.log.Warn("current answer not found", slog.String("err", err.Error()))
			return fmt.Errorf("%s: %w", op, ErrAnswerNotFound)
		}

		log.Error("failed to get current answer", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, ErrFailedToSaveInRedis)
	}

	pageAttemptToRedis := &redis.SavePageAttempt{
		LessonAttemptID: updPAttempt.LessonAttemptID,
		PageAttemptID:   updPAttempt.QPAttemptID,
		UserAnswer:      updPAttempt.UserAnswer,
	}

	if answer == updPAttempt.UserAnswer {
		pageAttemptToRedis.IsCorrect = true
	}

	// Save to Redis
	if err := ah.attemptRedisStore.SavePageAttempt(ctx, pageAttemptToRedis); err != nil {
		log.Error("failed to save page attempt in redis", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, ErrFailedToSaveInRedis)
	}

	return nil
}

func (ah *AttemptHandlers) CompleteLesson(ctx context.Context, req *attempts.CompleteLessonRequest) (*attempts.CompleteLessonResp, error) {
	const op = "attempts.CompleteLesson"

	log := ah.log.With(
		slog.String("op", op),
		slog.String("user_id", req.UserID),
		slog.Int64("lesson_attempt_id", req.LessonAttemptID),
	)

	// Validation
	err := ah.validator.Struct(req)
	if err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	// Get question page attempts from Redis
	pageAttemptsRds, err := ah.attemptRedisStore.GetPageAttempts(ctx, req.LessonAttemptID)
	if err != nil {
		log.Error("failed to get page attempts from redis", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	successCounter := 0
	pageAttemptsCount := len(pageAttemptsRds)
	currentTime := time.Now()
	// Update DB data
	for _, qPAttempt := range pageAttemptsRds {
		if err := ah.attemptSaver.UpdatePageAttempt(ctx, &attempts.UpdatePageAttempt{
			QPAttemptID:  qPAttempt.ID,
			UserAnswer:   qPAttempt.UserAnswer,
			Modified:     currentTime,
			IsSuccessful: qPAttempt.IsCorrect,
		}); err != nil {
			log.Error("failed to save attempt to DB", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if qPAttempt.IsCorrect {
			successCounter++
		}
	}

	updLessonAttempt := &attempts.UpdateLessonAttempt{
		UserID:          req.UserID,
		LessonAttemptID: req.LessonAttemptID,
		EndTime:         currentTime,
		IsComplete:      true,
	}

	// Colculate progress
	percentageScore := int64(0)

	if pageAttemptsCount == 0 {
		percentageScore = 100
		updLessonAttempt.PercentageScore = percentageScore
		updLessonAttempt.IsSuccessful = true
	} else {
		if pageAttemptsCount > 0 {
			percentageScore = int64(float64(successCounter) / float64(pageAttemptsCount) * 100)
		}
		updLessonAttempt.PercentageScore = percentageScore
		if percentageScore >= 75 {
			updLessonAttempt.IsSuccessful = true
		}
	}

	// Update lesson attempt
	id, err := ah.attemptSaver.UpdateLessonAttempt(ctx, updLessonAttempt)
	if err != nil {
		if errors.Is(err, storage.ErrLessonAttemtNotFound) {
			ah.log.Warn("lesson attempt not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrLessonAttemtNotFound)
		}

		log.Error("failed update lesson attempt", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrFailedToSaveInRedis)
	}

	return &attempts.CompleteLessonResp{
		ID:              id,
		IsSuccessful:    updLessonAttempt.IsSuccessful,
		PercentageScore: updLessonAttempt.PercentageScore,
	}, nil
}

func (ah *AttemptHandlers) GetLessonAttempts(ctx context.Context, inputParams *attempts.GetLessonAttempts) (*attempts.GetLessonAttemptsResp, error) {
	const op = "attempts.GetLessonAttempts"

	log := ah.log.With(
		slog.String("op", op),
		slog.String("user_id", inputParams.UserID),
	)

	// Validation
	params := attempts.GetLessonAttempts{
		UserID:   inputParams.UserID,
		LessonID: inputParams.LessonID,
		Limit:    inputParams.Limit,
		Offset:   inputParams.Offset,
	}
	params.SetDefaults()

	// Validation
	err := ah.validator.Struct(params)
	if err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("getting lesson attempts")

	attempts, err := ah.attemptProvider.GetLessonAttempts(ctx, &params)
	if err != nil {
		if errors.Is(err, storage.ErrLessonAttemtNotFound) {
			ah.log.Warn("lesson attempts not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrLessonAttemtNotFound)
		}

		log.Error("failed to get lesson attempts", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return attempts, nil
}

func (ah *AttemptHandlers) CheckPermissionForUser(ctx context.Context, userAtt *attempts.PermissionForUser) (bool, error) {
	const op = "attempts.CheckPermissionForUser"

	log := ah.log.With(
		slog.String("op", op),
		slog.String("user_id", userAtt.UserID),
		slog.Int64("lesson_attempt_id", userAtt.LessonAttemptID),
	)

	// Validation
	err := ah.validator.Struct(userAtt)
	if err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("checking permissions for user")

	perm, err := ah.attemptProvider.CheckPermissionForUser(ctx, userAtt)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, ErrPermissionsDenied)
	}

	return perm, nil
}

func (ah *AttemptHandlers) getExistingPageAttempts(ctx context.Context, lessonAttemptID int64, questionPage *attempts.GetQuestionPageAttempts, log *slog.Logger) ([]attempts.QuestionPageAttempt, error) {
	const op = "attempts.getExistingPageAttempts"

	// Get attempts from Redis
	pageAttemptsRds, err := ah.attemptRedisStore.GetPageAttempts(ctx, lessonAttemptID)
	if err != nil {
		log.Error("failed to get page attempts from redis", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrPageAttemtsNotFound)
	}

	if len(pageAttemptsRds) > 0 {
		return pageAttemptsRds, nil
	}

	// Fallback: Get attempts from DB
	pageAttemptsDB, err := ah.attemptProvider.GetLessonPagesAttempts(ctx, questionPage)
	if err != nil {
		if errors.Is(err, storage.ErrPageAttemtsNotFound) {
			log.Warn("lesson attempts not found", slog.String("err", err.Error()))
			return []attempts.QuestionPageAttempt{}, fmt.Errorf("%s: %w", op, ErrPageAttemtsNotFound)
		}
		log.Error("failed to get page attempts from DB", slog.String("err", err.Error()))
		return []attempts.QuestionPageAttempt{}, fmt.Errorf("%s: %w", op, err)
	}

	// Save attempts to Redis
	for _, attempt := range pageAttemptsDB {
		if err := ah.attemptRedisStore.SavePageAttempt(ctx, &redis.SavePageAttempt{
			LessonAttemptID: attempt.LessonAttemptID,
			PageAttemptID:   attempt.ID,
			UserAnswer:      "",
			IsCorrect:       false,
		}); err != nil {
			log.Error("failed to save page attempts in redis", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrFailedToSaveInRedis)
		}
	}

	return pageAttemptsDB, nil
}

func (ah *AttemptHandlers) createLessonAttemptAndPages(ctx context.Context, questionPage *attempts.GetQuestionPageAttempts, log *slog.Logger) ([]attempts.QuestionPageAttempt, error) {
	const op = "attempts.createLessonAttemptAndPages"

	// Create lesson attempt
	lAttemptID, err := ah.attemptSaver.CreateLessonAttempt(ctx, &attempts.CreateLessonAttempt{
		LessonID:  questionPage.LessonID,
		PlanID:    questionPage.PlanID,
		ChannelID: questionPage.ChannelID,
		UserID:    questionPage.UserID,
		StartTime: time.Now(),
	})
	if err != nil {
		log.Error("failed to save lesson attempt", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Get question pages
	qPages, err := ah.attemptProvider.GetQuestionPages(ctx, questionPage.LessonID)
	if err != nil {
		log.Error("failed to get question pages", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Create question page attempts
	return ah.createQuestionPageAttempts(ctx, lAttemptID, qPages, log)
}

func (ah *AttemptHandlers) createQuestionPageAttempts(ctx context.Context, lessonAttemptID int64, qPages []attempts.QuestionPage, log *slog.Logger) ([]attempts.QuestionPageAttempt, error) {
	const op = "attempts.createQuestionPageAttempts"

	var mappedAttempts []attempts.QuestionPageAttempt
	timeNow := time.Now()

	for _, qPage := range qPages {
		attempt, err := ah.attemptSaver.CreateQuestionPageAttempts(ctx, attempts.CreateQuestionPageAttemptNew{
			CreateAbstractPageAttempt: attempts.CreateAbstractPageAttempt{
				LessonAttemptID: lessonAttemptID,
				ContentType:     qPage.ContentType,
			},
			CreateAbstractQuestionAttempt: attempts.CreateAbstractQuestionAttempt{
				QuestionType: qPage.QuestionType,
			},
			CreateQuestionPageAttempt: attempts.CreateQuestionPageAttempt{
				PageID: qPage.QuestionPageID,
			},
			Modified: timeNow,
		})
		if err != nil {
			log.Error("failed to save question page attempt", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		// Save to Redis
		if err := ah.attemptRedisStore.SavePageAttempt(ctx, &redis.SavePageAttempt{
			LessonAttemptID: lessonAttemptID,
			PageAttemptID:   attempt.ID,
			UserAnswer:      "",
			IsCorrect:       false,
		}); err != nil {
			log.Error("failed to save page attempt in redis", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrFailedToSaveInRedis)
		}

		mappedAttempts = append(mappedAttempts, attempts.QuestionPageAttempt{
			ID:              attempt.ID,
			PageID:          attempt.PageID,
			LessonAttemptID: lessonAttemptID,
			UserAnswer:      "",
			IsCorrect:       false,
		})
	}

	return mappedAttempts, nil
}

// func (ah *AttemptHandlers) GetQuestionPageAttempts(ctx context.Context, questionPage *attempts.GetQuestionPageAttempts) ([]attempts.QuestionPageAttempt, error) {
// 	const op = "lesson.GetQuestionPageAttempts"

// 	log := ah.log.With(
// 		slog.String("op", op),
// 		slog.String("user_id", questionPage.UserID),
// 		slog.Int64("lesson_id", questionPage.LessonID),
// 	)

// 	// Validation
// 	err := ah.validator.Struct(questionPage)
// 	if err != nil {
// 		log.Warn("invalid parameters", slog.String("err", err.Error()))
// 		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
// 	}

// 	// Checks that lesson attempt is exist
// 	lessonAttemptID, err := ah.attemptProvider.CheckLessonAttempt(ctx, &attempts.GetQuestionPageAttempts{
// 		LessonID:  questionPage.LessonID,
// 		PlanID:    questionPage.PlanID,
// 		ChannelID: questionPage.ChannelID,
// 		UserID:    questionPage.UserID,
// 	})

// 	// If exist -> get page attempts and return it
// 	if lessonAttemptID != 0 {
// 		// Get from redis
// 		pageAttemptsRds, err := ah.attemptRedisStore.GetPageAttempts(ctx, lessonAttemptID)
// 		if err != nil {
// 			log.Error("failed to get page attempts from redis", slog.String("err", err.Error()))
// 			return nil, fmt.Errorf("%s: %w", op, err)
// 		}

// 		if len(pageAttemptsRds) == 0 {
// 			// Get from DB
// 			pageAttemptsDB, err := ah.attemptProvider.GetLessonPagesAttempts(ctx, questionPage)
// 			if err != nil {
// 				if errors.Is(err, storage.ErrPageAttemtsNotFound) {
// 					ah.log.Warn("lesson attemts not found", slog.String("err", err.Error()))
// 					return nil, fmt.Errorf("%s: %w", op, ErrPageAttemtsNotFound)
// 				}

// 				log.Error("failed to get page attempts", slog.String("err", err.Error()))
// 				return nil, fmt.Errorf("%s: %w", op, err)
// 			}

// 			// Save to Redis
// 			for _, attempt := range pageAttemptsDB {
// 				err := ah.attemptRedisStore.SavePageAttempt(ctx, &redis.SavePageAttempt{
// 					LessonAttemptID: attempt.LessonAttemptID,
// 					PageAttemptID:   attempt.ID,
// 					UserAnswer:      "",
// 					IsCorrect:       false,
// 				})
// 				if err != nil {
// 					log.Error("failed to save in redis", slog.String("err", err.Error()))
// 					return nil, fmt.Errorf("%s: %w", op, ErrFailedToSaveInRedis)
// 				}
// 			}

// 			return pageAttemptsDB, nil
// 		}

// 		return pageAttemptsRds, nil
// 	}

// 	lAttemptID, err := ah.attemptSaver.CreateLessonAttempt(ctx, &attempts.CreateLessonAttempt{
// 		LessonID:  questionPage.LessonID,
// 		PlanID:    questionPage.PlanID,
// 		ChannelID: questionPage.ChannelID,
// 		UserID:    questionPage.UserID,
// 		StartTime: time.Now(),
// 	})
// 	if err != nil {
// 		if errors.Is(err, storage.ErrInvalidCredentials) {
// 			ah.log.Warn("invalid arguments", slog.String("err", err.Error()))
// 			return nil, fmt.Errorf("%s: %w", op, err)
// 		}

// 		log.Error("failed to save lesson attempt", slog.String("err", err.Error()))
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	var qPages []attempts.QuestionPage
// 	qPages, err = ah.attemptProvider.GetQuestionPages(ctx, questionPage.LessonID)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrPageNotFound) || errors.Is(err, storage.ErrScanFailed) {
// 			ah.log.Warn("question pages not found", slog.String("err", err.Error()))
// 			return nil, fmt.Errorf("%s: %w", op, storage.ErrPageNotFound)
// 		}

// 		log.Error("failed to get question pages", slog.String("err", err.Error()))
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	var mappedQPageAttempts []attempts.QuestionPageAttempt
// 	timeNow := time.Now()
// 	for _, qPage := range qPages {
// 		qpAttempt, err := ah.attemptSaver.CreateQuestionPageAttempts(
// 			ctx,
// 			attempts.CreateQuestionPageAttemptNew{
// 				CreateAbstractPageAttempt: attempts.CreateAbstractPageAttempt{
// 					LessonAttemptID: lAttemptID,
// 					ContentType:     qPage.ContentType,
// 				},
// 				CreateAbstractQuestionAttempt: attempts.CreateAbstractQuestionAttempt{
// 					QuestionType: qPage.QuestionType,
// 				},
// 				CreateQuestionPageAttempt: attempts.CreateQuestionPageAttempt{
// 					PageID: qPage.QuestionPageID,
// 				},
// 				Modified: timeNow,
// 			},
// 		)
// 		if err != nil {
// 			if errors.Is(err, storage.ErrInvalidCredentials) {
// 				ah.log.Warn("invalid arguments", slog.String("err", err.Error()))
// 				return nil, fmt.Errorf("%s: %w", op, err)
// 			}

// 			log.Error("failed to save attempt", slog.String("err", err.Error()))
// 			return nil, fmt.Errorf("%s: %w", op, err)
// 		}

// 		err = ah.attemptRedisStore.SavePageAttempt(ctx, &redis.SavePageAttempt{
// 			LessonAttemptID: lAttemptID,
// 			PageAttemptID:   qpAttempt.ID,
// 			UserAnswer:      "",
// 			IsCorrect:       false,
// 		})
// 		if err != nil {
// 			log.Error("failed to save in redis", slog.String("err", err.Error()))
// 			return nil, fmt.Errorf("%s: %w", op, ErrFailedToSaveInRedis)
// 		}

// 		mappedQPageAttempts = append(mappedQPageAttempts, attempts.QuestionPageAttempt{
// 			ID:              qpAttempt.ID,
// 			PageID:          qpAttempt.PageID,
// 			LessonAttemptID: lAttemptID,
// 		})
// 	}

// 	return mappedQPageAttempts, nil
// }
