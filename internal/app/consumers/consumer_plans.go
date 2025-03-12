package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/DimTur/lp_learning_platform/internal/services/plan"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/plans"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchange   = "share"
	routingKey = "notification_to_auth"
)

type PlanStorage interface {
	plan.PlanSaver
	plan.PlanProvider
}

type RabbitMQQueues interface {
	plan.RabbitMQQueues
}

type ConsumerSharedPlans struct {
	msgQueue       MessageQueue
	planStorage    PlanStorage
	rabbitMQQueues RabbitMQQueues
	logger         *slog.Logger
}

func NewConsumePlan(
	msgQueue MessageQueue,
	planStorage PlanStorage,
	rabbitMQQueues RabbitMQQueues,
	logger *slog.Logger,
) *ConsumerSharedPlans {
	return &ConsumerSharedPlans{
		msgQueue:       msgQueue,
		planStorage:    planStorage,
		rabbitMQQueues: rabbitMQQueues,
		logger:         logger,
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

	// Decoding JSON message
	var message plans.SharePlanForUsers
	if err := json.Unmarshal(del.Body, &message); err != nil {
		c.logger.Error("failed to unmarshal message to SharePlanForUsers", slog.Any("err", err))
		return err
	}

	if err := c.planStorage.BatchSharePlansWithUsers(ctx, &message); err != nil {
		log.Error("failed to share plan", slog.String("err", err.Error()))
		log.Info(
			"failed to sharing",
			slog.Int64("splan", message.PlanID),
			slog.Any("with_users", message.UserIDs),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := c.rabbitMQQueues.Publish(ctx, exchange, routingKey, del.Body); err != nil {
		log.Error("err send sharing plan to exchange", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info(
		"successfully",
		slog.Int64("splan", message.PlanID),
		slog.Any("with_users", message.UserIDs),
	)

	return nil
}
