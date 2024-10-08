package attempt

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/DimTur/lp_learning_platform/internal/services/storage"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/attempts"
	"github.com/go-playground/validator/v10"
)

type AttemptSaver interface {
	CreateLessonAttempt(ctx context.Context, lAttempt attempts.CreateLessonAttempt) (int64, error)
	CreateAbstractPageAttempt(ctx context.Context, pAttempt attempts.CreateAbstractPageAttempt) (int64, error)
	CreateAbstractQuestionAttempt(ctx context.Context, qAttempt attempts.CreateAbstractQuestionAttempt) (int64, error)
	CreateQuestionAttempt(ctx context.Context, qPageAttempt attempts.CreateQuestionPageAttempt) error
}

type AttemptProvider interface {
	GetQuestionPages(ctx context.Context, lessonID int64) ([]attempts.QuestionPage, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAttemptID   = errors.New("invalid attempt id")
	ErrAttemptExitsts     = errors.New("attempt already exists")
	ErrAttemptNotFound    = errors.New("attempt not found")
	ErrFailedToCreate     = errors.New("attempt creation failed")
)

type AttemptHandlers struct {
	log             *slog.Logger
	validator       *validator.Validate
	attemptSaver    AttemptSaver
	attemptProvider AttemptProvider
}

func New(
	log *slog.Logger,
	validator *validator.Validate,
	attemptSaver AttemptSaver,
	attemptProvider AttemptProvider,
) *AttemptHandlers {
	return &AttemptHandlers{
		log:             log,
		validator:       validator,
		attemptSaver:    attemptSaver,
		attemptProvider: attemptProvider,
	}
}

// CreateAttempt creates new attempt of the lesson in the system and returns attempt ID.
func (ah *AttemptHandlers) CreateAttempt(ctx context.Context, attempt attempts.CreateLessonAttempt) (int64, error) {
	const op = "lesson.CreateAttempt"

	log := ah.log.With(
		slog.String("op", op),
		slog.Int64("lesson id", attempt.LessonID),
		slog.Int64("user id", attempt.UserID),
	)

	// Validation
	err := ah.validator.Struct(attempt)
	if err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("creating attempt")

	lAttemptID, err := ah.attemptSaver.CreateLessonAttempt(ctx, attempt)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidCredentials) {
			ah.log.Warn("invalid arguments", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to save lesson attempt", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var qPages []attempts.QuestionPage
	qPages, err = ah.attemptProvider.GetQuestionPages(ctx, attempt.LessonID)
	if err != nil {
		if errors.Is(err, storage.ErrPageNotFound) || errors.Is(err, storage.ErrScanFailed) {
			ah.log.Warn("question pages not found", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, ErrFailedToCreate)
		}

		log.Error("failed to get question pages", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	for _, qPage := range qPages {
		abPageID, err := ah.attemptSaver.CreateAbstractPageAttempt(ctx, attempts.CreateAbstractPageAttempt{
			LessonAttemptID: lAttemptID,
			ContentType:     qPage.ContentType,
		})
		if err != nil {
			if errors.Is(err, storage.ErrInvalidCredentials) {
				ah.log.Warn("invalid arguments", slog.String("err", err.Error()))
				return 0, fmt.Errorf("%s: %w", op, err)
			}

			log.Error("failed to save abstract page attempt", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		pageAttemptID, err := ah.attemptSaver.CreateAbstractQuestionAttempt(ctx, attempts.CreateAbstractQuestionAttempt{
			QuestionType:  qPage.QuestionType,
			PageAttemptID: abPageID,
		})
		if err != nil {
			if errors.Is(err, storage.ErrInvalidCredentials) {
				ah.log.Warn("invalid arguments", slog.String("err", err.Error()))
				return 0, fmt.Errorf("%s: %w", op, err)
			}

			log.Error("failed to save abstract question attempt", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		err = ah.attemptSaver.CreateQuestionAttempt(ctx, attempts.CreateQuestionPageAttempt{
			PageID:        qPage.QuestionPageID,
			PageAttemptID: pageAttemptID,
		})
		if err != nil {
			if errors.Is(err, storage.ErrInvalidCredentials) {
				ah.log.Warn("invalid arguments", slog.String("err", err.Error()))
				return 0, fmt.Errorf("%s: %w", op, err)
			}

			log.Error("failed to save question page attempt", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}

	return lAttemptID, nil
}
