package lesson

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/DimTur/lp_learning_platform/internal/services/storage"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/lessons"
	"github.com/go-playground/validator/v10"
)

type LessonSaver interface {
	CreateLesson(ctx context.Context, lesson *lessons.CreateLesson) (int64, error)
	UpdateLesson(ctx context.Context, updLesson *lessons.UpdateLessonRequest) (int64, error)
}

type LessonProvider interface {
	GetLessonByID(ctx context.Context, lessonPlan *lessons.GetLesson) (lessons.Lesson, error)
	GetLessons(ctx context.Context, inputParams *lessons.GetLessons) ([]lessons.Lesson, error)
}
type LessonDel interface {
	DeleteLesson(ctx context.Context, lessonP *lessons.DeleteLesson) error
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidLessonID    = errors.New("invalid lesson id")
	ErrLessonExitsts      = errors.New("lesson already exists")
	ErrLessonNotFound     = errors.New("lesson not found")
)

type LessonHandlers struct {
	log            *slog.Logger
	validator      *validator.Validate
	lessonSaver    LessonSaver
	lessonProvider LessonProvider
	lessonDel      LessonDel
}

func New(
	log *slog.Logger,
	validator *validator.Validate,
	lessonSaver LessonSaver,
	lessonProvider LessonProvider,
	lessonDel LessonDel,
) *LessonHandlers {
	return &LessonHandlers{
		log:            log,
		validator:      validator,
		lessonSaver:    lessonSaver,
		lessonProvider: lessonProvider,
		lessonDel:      lessonDel,
	}
}

// CreateLesson creats new lesson in the system and returns lesson ID.
func (lh *LessonHandlers) CreateLesson(ctx context.Context, lesson *lessons.CreateLesson) (int64, error) {
	const op = "lesson.CreateLesson"

	log := lh.log.With(
		slog.String("op", op),
		slog.String("lesson_name", lesson.Name),
	)

	// Validation
	err := lh.validator.Struct(lesson)
	if err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	now := time.Now()
	lesson.CreatedAt = now
	lesson.Modified = now

	log.Info("creating lesson")

	id, err := lh.lessonSaver.CreateLesson(ctx, lesson)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidCredentials) {
			lh.log.Warn("invalid arguments", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to save lesson", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// GetLesson gets lesson by ID and returns it.
func (lh *LessonHandlers) GetLesson(ctx context.Context, lessonPlan *lessons.GetLesson) (*lessons.Lesson, error) {
	const op = "lessons.GetLesson"

	log := lh.log.With(
		slog.String("op", op),
		slog.Int64("lesson_id", lessonPlan.LessonID),
		slog.Int64("plan_id", lessonPlan.PlanID),
	)

	log.Info("getting lesson")

	lesson, err := lh.lessonProvider.GetLessonByID(ctx, lessonPlan)
	if err != nil {
		if errors.Is(err, storage.ErrLessonNotFound) {
			lh.log.Warn("lesson not found", slog.String("err", err.Error()))
			return nil, ErrLessonNotFound
		}

		log.Error("failed to get lesson", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &lesson, nil
}

// GetLessons gets lessons and returns them.
func (lh *LessonHandlers) GetLessons(ctx context.Context, inputParams *lessons.GetLessons) ([]lessons.Lesson, error) {
	const op = "lessons.GetLessons"

	log := lh.log.With(
		slog.String("op", op),
		slog.Int64("getting lessons included in plan with id", inputParams.PlanID),
	)

	log.Info("getting lessons")

	// Validation
	params := lessons.GetLessons{
		PlanID: inputParams.PlanID,
		Limit:  inputParams.Limit,
		Offset: inputParams.Offset,
	}
	params.SetDefaults()

	if err := lh.validator.Struct(params); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	lessons, err := lh.lessonProvider.GetLessons(ctx, inputParams)
	if err != nil {
		if errors.Is(err, storage.ErrLessonNotFound) {
			lh.log.Warn("lessons not found", slog.String("err", err.Error()))
			return lessons, fmt.Errorf("%s: %w", op, ErrLessonNotFound)
		}

		log.Error("failed to get lessons", slog.String("err", err.Error()))
		return lessons, fmt.Errorf("%s: %w", op, err)
	}

	return lessons, nil
}

// UpdateLesson performs a partial update
func (lh *LessonHandlers) UpdateLesson(ctx context.Context, updLesson *lessons.UpdateLessonRequest) (int64, error) {
	const op = "lessons.UpdateLesson"

	log := lh.log.With(
		slog.String("op", op),
		slog.Int64("lesson_id: ", updLesson.LessonID),
		slog.Int64("plan_id: ", updLesson.PlanID),
	)

	log.Info("updating lesson")

	// Validation
	err := lh.validator.Struct(updLesson)
	if err != nil {
		log.Warn("validation failed", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	id, err := lh.lessonSaver.UpdateLesson(ctx, updLesson)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrLessonNotFound):
			lh.log.Warn("lesson not found", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, ErrLessonNotFound)
		case errors.Is(err, storage.ErrInvalidCredentials):
			lh.log.Warn("invalid credentials", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to update lesson", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}
	log.Info("lesson updated with ", slog.Int64("lesson_id", id))

	return id, nil
}

// DeleteLesson
func (lh *LessonHandlers) DeleteLesson(ctx context.Context, lessonP *lessons.DeleteLesson) error {
	const op = "lessons.DeleteLesson"

	log := lh.log.With(
		slog.String("op", op),
		slog.Int64("lesson_id", lessonP.LessonID),
		slog.Int64("plan_id", lessonP.PlanID),
	)

	log.Info("deleting lesson with: ", slog.Int64("lesson_id", lessonP.LessonID))

	err := lh.lessonDel.DeleteLesson(ctx, lessonP)
	if err != nil {
		if errors.Is(err, storage.ErrLessonNotFound) {
			lh.log.Warn("lesson not found", slog.String("err", err.Error()))
			return fmt.Errorf("%s: %w", op, ErrLessonNotFound)
		}

		log.Error("failed to delete lesson", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
