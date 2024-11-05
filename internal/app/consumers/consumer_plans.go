package consumers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/DimTur/lp_learning_platform/internal/services/channel"
	"github.com/DimTur/lp_learning_platform/internal/services/plan"
	"github.com/DimTur/lp_learning_platform/internal/services/storage"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/plans"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PlanStorage interface {
	plan.PlanSaver
}

type ConsumerSharedPlans struct {
	msgQueue    MessageQueue
	planStorage PlanStorage
	logger      *slog.Logger
}

func NewConsumePlan(
	msgQueue MessageQueue,
	planStorage PlanStorage,
	logger *slog.Logger,
) *ConsumerSharedPlans {
	return &ConsumerSharedPlans{
		msgQueue:    msgQueue,
		planStorage: planStorage,
		logger:      logger,
	}
}

func (c *ConsumerSharedPlans) Start(ctx context.Context,
	queueName, consumer string,
	autoAck, exclusive, noLocal, noWait bool,
	args map[string]interface{},
) error {
	const op = "ConsumerSharedPlans.Start"

	log := c.logger.With(slog.String("op", op))
	log.Info("Starting to consume shared plans messages")

	return c.msgQueue.Consume(
		ctx,
		queueName,
		consumer,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		args,
		c.handleMessage)
}

func (c *ConsumerSharedPlans) handleMessage(ctx context.Context, msg interface{}) error {
	const op = "consumer_channels.handleMessage"

	log := c.logger.With(
		slog.String("op", op),
	)

	del, ok := msg.(amqp.Delivery)
	if !ok {
		c.logger.Error("failed to cast message to amqp.Delivery")
		return nil // Return nil to avoid calling Nack/Ack
	}

	var message plans.SharePlanForUsers
	// Decoding JSON message
	if err := json.Unmarshal(del.Body, &message); err != nil {
		c.logger.Error("failed to unmarshal message to SharePlanForUsers", slog.Any("err", err))
		return err
	}

	for _, userID := range message.UsersIDs {
		share := &plans.DBSharePlanForUser{
			PlanID:    message.PlanID,
			UserID:    userID,
			CreatedBy: message.CreatedBy,
			CreatedAt: time.Now(),
		}

		if err := c.planStorage.SharePlanWithUser(ctx, *share); err != nil {
			if errors.Is(err, storage.ErrInvalidCredentials) {
				log.Warn("invalid arguments", slog.String("err", err.Error()))
				log.Info(
					"failed to sharing",
					slog.Int64("splan", share.PlanID),
					slog.String("with user", userID),
				)
				return fmt.Errorf("%s: %w", op, channel.ErrInvalidCredentials)
			}

			log.Error("failed to share plan", slog.String("err", err.Error()))
			log.Info(
				"failed to sharing",
				slog.Int64("splan", share.PlanID),
				slog.String("with user", userID),
			)
			return fmt.Errorf("%s: %w", op, err)
		}
		log.Info(
			"successfully",
			slog.Int64("shared plan", share.PlanID),
			slog.String("with user", userID),
		)
	}

	return nil
}
