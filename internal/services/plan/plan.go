package plan

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/DimTur/lp_learning_platform/internal/domain/models"
	"github.com/DimTur/lp_learning_platform/internal/services/storage"
	"github.com/DimTur/lp_learning_platform/internal/utils"
	"github.com/go-playground/validator/v10"
)

type PlanSaver interface {
	CreatePlan(ctx context.Context, plan models.CreatePlan) (id int64, err error)
	UpdatePlan(ctx context.Context, updPlan models.UpdatePlanRequest) (id int64, err error)
}

type PlanProvider interface {
	GetPlanByID(ctx context.Context, planID int64) (plan models.Plan, err error)
	GetPlans(ctx context.Context, channel_id int64, limit, offset int64) (plans []models.Plan, err error)
}
type PlanDel interface {
	DeletePlan(ctx context.Context, id int64) (err error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPlanID      = errors.New("invalid plan id")
	ErrPlanExitsts        = errors.New("plan already exists")
	ErrPlanNotFound       = errors.New("plan not found")
)

type PlanHandlers struct {
	log          *slog.Logger
	validator    *validator.Validate
	planSaver    PlanSaver
	planProvider PlanProvider
	planDel      PlanDel
}

func New(
	log *slog.Logger,
	validator *validator.Validate,
	planSaver PlanSaver,
	planProvider PlanProvider,
	planDel PlanDel,
) *PlanHandlers {
	return &PlanHandlers{
		log:          log,
		validator:    validator,
		planSaver:    planSaver,
		planProvider: planProvider,
		planDel:      planDel,
	}
}

// CreatePlan creats new plan in the system and returns plan ID.
func (ph *PlanHandlers) CreatePlan(ctx context.Context, plan models.CreatePlan) (int64, error) {
	const op = "plan.CreatePlan"

	log := ph.log.With(
		slog.String("op", op),
		slog.String("name", plan.Name),
	)

	// Validation
	err := ph.validator.Struct(plan)
	if err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	now := time.Now()
	plan.CreatedAt = now
	plan.Modified = now
	plan.IsPublished = false
	plan.Public = false

	log.Info("creating plan")

	id, err := ph.planSaver.CreatePlan(ctx, plan)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidCredentials) {
			ph.log.Warn("invalid arguments", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to save plan", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// GetPlan gets plan by ID and returns it.
func (ph *PlanHandlers) GetPlan(ctx context.Context, planID int64) (models.Plan, error) {
	const op = "plans.GetPlan"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("planID", planID),
	)

	log.Info("getting plan")

	var plan models.Plan
	plan, err := ph.planProvider.GetPlanByID(ctx, planID)
	if err != nil {
		if errors.Is(err, storage.ErrPlanNotFound) {
			ph.log.Warn("plan not found", slog.String("err", err.Error()))
			return plan, ErrPlanNotFound
		}

		log.Error("failed to get plan", slog.String("err", err.Error()))
		return plan, fmt.Errorf("%s: %w", op, err)
	}

	return plan, nil
}

// GetPlans gets plans and returns them.
func (ph *PlanHandlers) GetPlans(ctx context.Context, channel_id int64, limit, offset int64) ([]models.Plan, error) {
	const op = "plans.GetPlans"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("getting plans included in channel with id", channel_id),
	)

	log.Info("getting plans")

	// Validation
	params := utils.PaginationQueryParams{
		Limit:  limit,
		Offset: offset,
	}

	if err := ph.validator.Struct(params); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	var plans []models.Plan
	plans, err := ph.planProvider.GetPlans(ctx, channel_id, limit, offset)
	if err != nil {
		if errors.Is(err, storage.ErrPlanNotFound) {
			ph.log.Warn("plans not found", slog.String("err", err.Error()))
			return plans, fmt.Errorf("%s: %w", op, ErrPlanNotFound)
		}

		log.Error("failed to get plans", slog.String("err", err.Error()))
		return plans, fmt.Errorf("%s: %w", op, err)
	}

	return plans, nil
}

// UpdatePlan performs a partial update
func (ph *PlanHandlers) UpdatePlan(ctx context.Context, updPlan models.UpdatePlanRequest) (int64, error) {
	const op = "plans.UpdatePlan"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("updating plan with id: ", updPlan.ID),
	)

	log.Info("updating plan")

	// Validation
	err := ph.validator.Struct(updPlan)
	if err != nil {
		log.Warn("validation failed", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	id, err := ph.planSaver.UpdatePlan(ctx, updPlan)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidCredentials) {
			ph.log.Warn("invalid credentials", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to update plan", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("plan updated with ", slog.Int64("planID", id))

	return id, nil
}

// DeletePlan
func (ph *PlanHandlers) DeletePlan(ctx context.Context, planID int64) error {
	const op = "plans.DeletePlan"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("plan id", planID),
	)

	log.Info("deleting plan with: ", slog.Int64("planID", planID))

	err := ph.planDel.DeletePlan(ctx, planID)
	if err != nil {
		if errors.Is(err, storage.ErrPlanNotFound) {
			ph.log.Warn("plan not found", slog.String("err", err.Error()))
			return fmt.Errorf("%s: %w", op, ErrPlanNotFound)
		}

		log.Error("failed to delete plan", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
