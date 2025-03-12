package consumers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/DimTur/lp_learning_platform/internal/services/channel"
	"github.com/DimTur/lp_learning_platform/internal/services/storage"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/channels"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueue interface {
	Consume(
		ctx context.Context,
		queueName, consumer string,
		autoAck, exclusive, noLocal, noWait bool,
		args map[string]interface{},
		handle func(ctx context.Context, msg interface{}) error,
	) error
}

type ChannelStorage interface {
	channel.ChannelSaver
}

type ConsumerSharedChannels struct {
	msgQueue       MessageQueue
	channelStorage ChannelStorage
	logger         *slog.Logger
}

func NewConsumeChannel(
	msgQueue MessageQueue,
	channelStorage ChannelStorage,
	logger *slog.Logger,
) *ConsumerSharedChannels {
	return &ConsumerSharedChannels{
		msgQueue:       msgQueue,
		channelStorage: channelStorage,
		logger:         logger,
	}
}

func (c *ConsumerSharedChannels) Start(ctx context.Context,
	queueName, consumer string,
	autoAck, exclusive, noLocal, noWait bool,
	args map[string]interface{},
) error {
	const op = "ConsumerSharedChannels.Start"

	log := c.logger.With(slog.String("op", op))
	log.Info("Starting to consume shared channels messages")

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

func (c *ConsumerSharedChannels) handleMessage(ctx context.Context, msg interface{}) error {
	const op = "consumer_channels.handleMessage"

	del, ok := msg.(amqp.Delivery)
	if !ok {
		c.logger.Error("failed to cast message to amqp.Delivery")
		return nil // Return nil to avoid calling Nack/Ack
	}

	var message channels.ShareChannelToGroup
	// Decoding JSON message
	if err := json.Unmarshal(del.Body, &message); err != nil {
		c.logger.Error("failed to unmarshal message to ShareChannelToGroup", slog.Any("err", err))
		return err
	}

	for _, groupID := range message.LGroupIDs {
		share := &channels.DBShareChannelToGroup{
			ChannelID: message.ChannelID,
			LGroupID:  groupID,
			CreatedBy: message.CreatedBy,
			CreatedAt: time.Now(),
		}

		if err := c.channelStorage.ShareChannelToGroup(ctx, *share); err != nil {
			if errors.Is(err, storage.ErrInvalidCredentials) {
				c.logger.Warn("invalid arguments", slog.String("err", err.Error()))
				c.logger.Info(
					"failed to sharing",
					slog.Int64("channel", share.ChannelID),
					slog.String("with group", groupID),
				)
				return fmt.Errorf("%s: %w", op, channel.ErrInvalidCredentials)
			}

			c.logger.Error("failed to share channel", slog.String("err", err.Error()))
			c.logger.Info(
				"failed to sharing",
				slog.Int64("channel", share.ChannelID),
				slog.String("with group", groupID),
			)
			return fmt.Errorf("%s: %w", op, err)
		}
		c.logger.Info(
			"successfully",
			slog.Int64("shared channel", share.ChannelID),
			slog.String("with group", groupID),
		)
	}

	return nil
}
