package question

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/DimTur/lp_learning_platform/internal/services/storage"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/questions"
	"github.com/go-playground/validator/v10"
)

type QuestionPageSaver interface {
	CreateQuestionPage(ctx context.Context, questionPage questions.CreateQuestionPage) (id int64, err error)
	UpdateQuestionPage(ctx context.Context, updPage questions.UpdateQuestionPage) (id int64, err error)
}

type QuestionPageProvider interface {
	GetQuestionPageByID(ctx context.Context, pageID int64) (questionPage questions.QuestionPage, err error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPageID      = errors.New("invalid page id")
	ErrPageExitsts        = errors.New("page already exists")
	ErrPageNotFound       = errors.New("page not found")
)

type QuestionPageHandlers struct {
	log                  *slog.Logger
	validator            *validator.Validate
	questionPageSaver    QuestionPageSaver
	questionPageProvider QuestionPageProvider
}

func New(
	log *slog.Logger,
	validator *validator.Validate,
	questionPageSaver QuestionPageSaver,
	questionPageProvider QuestionPageProvider,
) *QuestionPageHandlers {
	return &QuestionPageHandlers{
		log:                  log,
		validator:            validator,
		questionPageSaver:    questionPageSaver,
		questionPageProvider: questionPageProvider,
	}
}

// CreateQuestionPage creates new question page in the system and returns page ID.
func (qph QuestionPageHandlers) CreateQuestionPage(ctx context.Context, questionPage questions.CreateQuestionPage) (int64, error) {
	const op = "question.CreateQuestionPage"

	log := qph.log.With(
		slog.String("op", op),
		slog.String("page with type", questionPage.ContentType),
	)

	// Validation
	err := qph.validator.Struct(questionPage)
	if err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("creating question page")

	id, err := qph.questionPageSaver.CreateQuestionPage(ctx, questionPage)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidCredentials) {
			qph.log.Warn("invalid arguments", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to save question page", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("question page created with %s", slog.Int64("id:", id))

	return id, nil
}

// GetQuestionPageByID gets question page by ID and returns it.
func (qph QuestionPageHandlers) GetQuestionPageByID(ctx context.Context, pageID int64) (questions.QuestionPage, error) {
	const op = "question.GetQuestionPageByID"

	log := qph.log.With(
		slog.String("op", op),
		slog.Int64("pageID", pageID),
	)

	log.Info("getting question page")

	var questionPage questions.QuestionPage
	questionPage, err := qph.questionPageProvider.GetQuestionPageByID(ctx, pageID)
	if err != nil {
		if errors.Is(err, storage.ErrPageNotFound) {
			qph.log.Warn("question page not found", slog.String("err", err.Error()))
			return questionPage, ErrPageNotFound
		}

		log.Error("failed to get question page", slog.String("err", err.Error()))
		return questionPage, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("question page received with %s", slog.Int64("id:", pageID))

	return questionPage, nil
}

// UpdateQuestionPage performs a partial update
func (qph QuestionPageHandlers) UpdateQuestionPage(ctx context.Context, updPage questions.UpdateQuestionPage) (int64, error) {
	const op = "question.UpdateQuestionPage"

	log := qph.log.With(
		slog.String("op", op),
		slog.Int64("updating question page with id: ", updPage.ID),
	)

	log.Info("updating question page")

	// Validation
	err := qph.validator.Struct(updPage)
	if err != nil {
		log.Warn("validation failed", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	id, err := qph.questionPageSaver.UpdateQuestionPage(ctx, updPage)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidCredentials) {
			qph.log.Warn("invalid credentials", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to update question page", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("question page updated with", slog.Int64("id:", id))

	return id, nil
}
