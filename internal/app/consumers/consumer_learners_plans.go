package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/DimTur/lp_learning_platform/internal/services/rabbitmq"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/plans"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangePlan   = "share"
	planRoutingKey = "plan"
)

type ConsumerSharedLearnersWithPlan struct {
	msgQueue       MessageQueue
	planStorage    PlanStorage
	rabbitMQQueues RabbitMQQueues
	logger         *slog.Logger
}

func NewConsumeSharedLearnersWithPlan(
	msgQueue MessageQueue,
	planStorage PlanStorage,
	rabbitMQQueues RabbitMQQueues,
	logger *slog.Logger,
) *ConsumerSharedLearnersWithPlan {
	return &ConsumerSharedLearnersWithPlan{
		msgQueue:       msgQueue,
		planStorage:    planStorage,
		rabbitMQQueues: rabbitMQQueues,
		logger:         logger,
	}
}

func (c *ConsumerSharedLearnersWithPlan) Start(ctx context.Context,
	queueName, consumer string,
	autoAck, exclusive, noLocal, noWait bool,
	args map[string]interface{},
) error {
	const op = "ConsumerSharedLearnersWithPlan.Start"

	log := c.logger.With(slog.String("op", op))
	log.Info("Starting to consume shared plans with learners messages")

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

func (c *ConsumerSharedLearnersWithPlan) handleMessage(ctx context.Context, msg interface{}) error {
	const op = "consumer_channels.handleMessage"

	log := c.logger.With(
		slog.String("op", op),
	)

	del, ok := msg.(amqp.Delivery)
	if !ok {
		c.logger.Error("failed to cast message to amqp.Delivery")
		return nil // Return nil to avoid calling Nack/Ack
	}

	// Decoding JSON message
	var message rabbitmq.Spfu
	if err := json.Unmarshal(del.Body, &message); err != nil {
		c.logger.Error("failed to unmarshal message to SharePlanForUsers", slog.Any("err", err))
		return err
	}

	// Get channels ids with plans ids
	planChannelIDs, err := c.planStorage.GetPlansForSharing(ctx, &plans.LearningGroup{
		LgID: message.LearningGroupID,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	now := time.Now()
	for key, value := range planChannelIDs {
		for _, plan := range value {
			newMsg := &plans.SharePlanForUsers{
				ChannelID: key,
				PlanID:    plan,
				UserIDs:   message.UserIDs,
				CreatedBy: message.CreatedBy,
				CreatedAt: now,
			}

			// Serialization and publication message
			msgBody, err := json.Marshal(newMsg)
			if err != nil {
				log.Error("failed to marshal newMsg request", slog.String("err", err.Error()))
				return fmt.Errorf("%s: %w", op, err)
			}

			if err = c.rabbitMQQueues.Publish(ctx, exchangePlan, planRoutingKey, msgBody); err != nil {
				log.Error("failed to publish newMsg request to plan queue", slog.String("err", err.Error()))
				return fmt.Errorf("%s: %w", op, err)
			}

			log.Info(
				"successfully",
				slog.Int64("sharing_plan_id", newMsg.PlanID),
				slog.Any("with_users", message.UserIDs),
			)
		}
	}

	return nil
}
