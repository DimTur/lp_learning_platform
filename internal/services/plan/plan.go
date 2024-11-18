package plan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/DimTur/lp_learning_platform/internal/services/storage"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/plans"
	"github.com/go-playground/validator/v10"
)

const (
	exchangePlan   = "share"
	queuePlan      = "plan"
	planRoutingKey = "plan"
)

type PlanSaver interface {
	CreatePlan(ctx context.Context, plan plans.CreatePlan) (int64, error)
	UpdatePlan(ctx context.Context, updPlan *plans.UpdatePlanRequest) (int64, error)
	SharePlanWithUser(ctx context.Context, s *plans.DBSharePlanForUser) error
}

type PlanProvider interface {
	GetPlanByID(ctx context.Context, planCh *plans.GetPlan) (plans.Plan, error)
	GetPlans(ctx context.Context, inputParams *plans.GetPlans) ([]plans.Plan, error)
	IsUserShareWithPlan(ctx context.Context, userPlan *plans.IsUserShareWithPlan) (bool, error)
	CanShare(ctx context.Context, cs *plans.DBCanShare) (bool, error)
}
type PlanDel interface {
	DeletePlan(ctx context.Context, planCh *plans.DeletePlan) error
}

type RabbitMQQueues interface {
	Publish(ctx context.Context, exchange, routingKey string, body []byte) error
	PublishToQueue(ctx context.Context, queueName string, body []byte) error
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPlanID      = errors.New("invalid plan id")
	ErrPlanExitsts        = errors.New("plan already exists")
	ErrPlanNotFound       = errors.New("plan not found")
)

type PlanHandlers struct {
	log            *slog.Logger
	validator      *validator.Validate
	planSaver      PlanSaver
	planProvider   PlanProvider
	planDel        PlanDel
	rabbitMQQueues RabbitMQQueues
}

func New(
	log *slog.Logger,
	validator *validator.Validate,
	planSaver PlanSaver,
	planProvider PlanProvider,
	planDel PlanDel,
	rabbitMQQueues RabbitMQQueues,
) *PlanHandlers {
	return &PlanHandlers{
		log:            log,
		validator:      validator,
		planSaver:      planSaver,
		planProvider:   planProvider,
		planDel:        planDel,
		rabbitMQQueues: rabbitMQQueues,
	}
}

// CreatePlan creats new plan in the system and returns plan ID.
func (ph *PlanHandlers) CreatePlan(ctx context.Context, plan plans.CreatePlan) (int64, error) {
	const op = "plan.CreatePlan"

	log := ph.log.With(
		slog.String("op", op),
		slog.String("name", plan.Name),
	)

	now := time.Now()
	plan.LastModifiedBy = plan.CreatedBy
	plan.CreatedAt = now
	plan.Modified = now
	plan.IsPublished = false
	plan.Public = false

	// Validation
	err := ph.validator.Struct(plan)
	if err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("creating plan")

	id, err := ph.planSaver.CreatePlan(ctx, plan)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidCredentials) {
			ph.log.Warn("invalid arguments", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to save plan", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	} else {
		s := &plans.SharePlanForUsers{
			PlanID:    id,
			UsersIDs:  []string{plan.CreatedBy},
			CreatedBy: plan.CreatedBy,
		}
		msgBody, err := json.Marshal(s)
		if err != nil {
			ph.log.Error("err to marshal shared msg", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		if err = ph.rabbitMQQueues.Publish(ctx, exchangePlan, planRoutingKey, msgBody); err != nil {
			ph.log.Error("err send sharing plan to exchange", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}

	return id, nil
}

// GetPlan gets plan by ID and returns it.
func (ph *PlanHandlers) GetPlan(ctx context.Context, planCh *plans.GetPlan) (plans.Plan, error) {
	const op = "plans.GetPlan"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("channel_id", planCh.ChannelID),
		slog.Int64("plan_id", planCh.PlanID),
	)

	log.Info("getting plan")

	plan, err := ph.planProvider.GetPlanByID(ctx, planCh)
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
func (ph *PlanHandlers) GetPlans(ctx context.Context, inputParams *plans.GetPlans) ([]plans.Plan, error) {
	const op = "plans.GetPlans"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("getting plans included in channel with id", inputParams.ChannelID),
	)

	log.Info("getting plans")

	// Validation
	params := plans.GetPlans{
		UserID:    inputParams.UserID,
		ChannelID: inputParams.ChannelID,
		Limit:     inputParams.Limit,
		Offset:    inputParams.Offset,
	}
	params.SetDefaults()

	if err := ph.validator.Struct(params); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	var plans []plans.Plan
	plans, err := ph.planProvider.GetPlans(ctx, &params)
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
func (ph *PlanHandlers) UpdatePlan(ctx context.Context, updPlan *plans.UpdatePlanRequest) (int64, error) {
	const op = "plans.UpdatePlan"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("channel_id", updPlan.ChannelID),
		slog.Int64("plan_id", updPlan.PlanID),
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
		switch {
		case errors.Is(err, storage.ErrPlanNotFound):
			ph.log.Warn("plan not found", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, ErrPlanNotFound)
		case errors.Is(err, storage.ErrInvalidCredentials):
			ph.log.Warn("invalid credentials", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		default:
			log.Error("failed to update plan", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}
	log.Info("plan updated with ", slog.Int64("planID", id))

	return id, nil
}

// DeletePlan
func (ph *PlanHandlers) DeletePlan(ctx context.Context, planCh *plans.DeletePlan) error {
	const op = "plans.DeletePlan"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("channel_id", planCh.ChannelID),
		slog.Int64("plan_id", planCh.PlanID),
	)

	log.Info("deleting plan with: ", slog.Int64("plan_id", planCh.PlanID))

	err := ph.planDel.DeletePlan(ctx, planCh)
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

// SharePlanWithUser sharing channel with lerning group
func (ph *PlanHandlers) SharePlanWithUser(ctx context.Context, s *plans.SharePlanForUsers) error {
	const op = "plan.SharePlanWithUser"

	log := ph.log.With(
		slog.String("op", op),
		slog.Int64("plan_id", s.PlanID),
		slog.String("created_by", s.CreatedBy),
	)

	// Validation
	err := ph.validator.Struct(s)
	if err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	canShare, err := ph.planProvider.CanShare(ctx, &plans.DBCanShare{
		ChannelID: s.ChannelID,
		PlanID:    s.PlanID,
	})
	if err != nil {
		ph.log.Error("invalid credentials", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	if !canShare {
		ph.log.Error("can't sharing")
		return fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	msgBody, err := json.Marshal(s)
	if err != nil {
		ph.log.Error("err to marshal shared msg", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = ph.rabbitMQQueues.Publish(ctx, exchangePlan, planRoutingKey, msgBody); err != nil {
		ph.log.Error("err send sharing plan to exchange", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("plan sent to share with users")

	return nil
}

// IsUserShareWithPlan checks that plan share with user
func (ph *PlanHandlers) IsUserShareWithPlan(ctx context.Context, userPlan *plans.IsUserShareWithPlan) (bool, error) {
	const op = "plan.IsUserShareWithPlan"

	log := ph.log.With(
		slog.String("op", op),
		slog.String("user_id", userPlan.UserID),
		slog.Int64("plan_id", userPlan.PlanID),
	)

	// Validation
	err := ph.validator.Struct(userPlan)
	if err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	isShare, err := ph.planProvider.IsUserShareWithPlan(ctx, userPlan)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isShare, nil
}
