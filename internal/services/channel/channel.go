package channel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/DimTur/lp_learning_platform/internal/services/storage"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/channels"
	"github.com/DimTur/lp_learning_platform/internal/utils"
	"github.com/go-playground/validator/v10"
)

const (
	exchangeChannel   = "share"
	queueChannel      = "channel"
	channelRoutingKey = "channel"
)

type ChannelSaver interface {
	CreateChannel(ctx context.Context, channel channels.CreateChannel) (int64, error)
	UpdateChannel(ctx context.Context, updChannel channels.UpdateChannelRequest) (int64, error)
	ShareChannelToGroup(ctx context.Context, s channels.DBShareChannelToGroup) error
}

type ChannelProvider interface {
	GetChannelByID(ctx context.Context, channelID int64) (channels.ChannelWithPlans, error)
	GetChannels(ctx context.Context, limit, offset int64) ([]channels.Channel, error)
}

type ChannelDel interface {
	DeleteChannel(ctx context.Context, channelID int64) error
}

type RabbitMQQueues interface {
	Publish(ctx context.Context, exchange, routingKey string, body []byte) error
	PublishToQueue(ctx context.Context, queueName string, body []byte) error
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidChannelID   = errors.New("invalid channel id")
	ErrChannelExitsts     = errors.New("channel already exists")
	ErrChannelNotFound    = errors.New("channel not found")
)

type ChannelHandlers struct {
	log             *slog.Logger
	validator       *validator.Validate
	channelSaver    ChannelSaver
	channelProvider ChannelProvider
	channelDel      ChannelDel
	rabbitMQQueues  RabbitMQQueues
}

func New(
	log *slog.Logger,
	validator *validator.Validate,
	channelSaver ChannelSaver,
	channelProvider ChannelProvider,
	channelDel ChannelDel,
	rabbitMQQueues RabbitMQQueues,
) *ChannelHandlers {
	return &ChannelHandlers{
		log:             log,
		validator:       validator,
		channelSaver:    channelSaver,
		channelProvider: channelProvider,
		channelDel:      channelDel,
		rabbitMQQueues:  rabbitMQQueues,
	}
}

// CreateChannel creats new channel in the system and returns channel ID.
func (chh *ChannelHandlers) CreateChannel(ctx context.Context, channel channels.CreateChannel) (int64, error) {
	const op = "channel.CreateChannel"

	log := chh.log.With(
		slog.String("op", op),
		slog.String("name", channel.Name),
	)

	// Validation
	err := chh.validator.Struct(channel)
	if err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	now := time.Now()
	channel.CreatedAt = now
	channel.Modified = now

	log.Info("creating channel")

	id, err := chh.channelSaver.CreateChannel(ctx, channel)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidCredentials) {
			chh.log.Warn("invalid arguments", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to save channel", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// GetChannelByID gets channel by ID and returns it.
func (chh *ChannelHandlers) GetChannel(ctx context.Context, channelID int64) (channels.ChannelWithPlans, error) {
	const op = "channel.GetChannelByID"

	log := chh.log.With(
		slog.String("op", op),
		slog.Int64("chanID", channelID),
	)

	log.Info("getting channel")

	var channel channels.ChannelWithPlans
	channel, err := chh.channelProvider.GetChannelByID(ctx, channelID)
	if err != nil {
		if errors.Is(err, storage.ErrChannelNotFound) {
			chh.log.Warn("channel not found", slog.String("err", err.Error()))
			return channel, ErrChannelNotFound
		}

		log.Error("failed to get channel", slog.String("err", err.Error()))
		return channel, fmt.Errorf("%s: %w", op, err)
	}

	return channel, nil
}

// GetChannels gets channels and returns them.
func (chh *ChannelHandlers) GetChannels(ctx context.Context, limit, offset int64) ([]channels.Channel, error) {
	const op = "channel.GetChannels"

	log := chh.log.With(
		slog.String("op", op),
	)

	log.Info("getting channels")

	// Validation
	params := utils.PaginationQueryParams{
		Limit:  limit,
		Offset: offset,
	}
	params.SetDefaults()

	if err := chh.validator.Struct(params); err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	var channels []channels.Channel
	channels, err := chh.channelProvider.GetChannels(ctx, params.Limit, params.Offset)
	if err != nil {
		if errors.Is(err, storage.ErrChannelNotFound) {
			chh.log.Warn("channels not found", slog.String("err", err.Error()))
			return channels, fmt.Errorf("%s: %w", op, ErrChannelNotFound)
		}

		log.Error("failed to get channels", slog.String("err", err.Error()))
		return channels, fmt.Errorf("%s: %w", op, err)
	}

	return channels, nil
}

// UpdateChannel performs a partial update
func (chh *ChannelHandlers) UpdateChannel(ctx context.Context, updChannel channels.UpdateChannelRequest) (int64, error) {
	const op = "channel.UpdateChannel"

	log := chh.log.With(
		slog.String("op", op),
	)

	log.Info("updating channel")

	// Validation
	err := chh.validator.Struct(updChannel)
	if err != nil {
		log.Warn("validation failed", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	id, err := chh.channelSaver.UpdateChannel(ctx, updChannel)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidCredentials) {
			chh.log.Warn("invalid credentials", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to update channel", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("channel updated with ", slog.Int64("channelID", id))

	return id, nil
}

// DeleteChannel
func (chh *ChannelHandlers) DeleteChannel(ctx context.Context, channelID int64) error {
	const op = "channel.DeleteChannel"

	log := chh.log.With(
		slog.String("op", op),
		slog.Int64("channel id", channelID),
	)

	log.Info("deleting channel with: ", slog.Int64("channelID", channelID))

	err := chh.channelDel.DeleteChannel(ctx, channelID)
	if err != nil {
		if errors.Is(err, storage.ErrChannelNotFound) {
			chh.log.Warn("channel not found", slog.String("err", err.Error()))
			return fmt.Errorf("%s: %w", op, ErrChannelNotFound)
		}

		log.Error("failed to delete channel", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// ShareChannelToGroup sharing channel with lerning group
func (chh *ChannelHandlers) ShareChannelToGroup(ctx context.Context, s channels.ShareChannelToGroup) error {
	const op = "channel.ShareChannelToGroup"

	log := chh.log.With(
		slog.String("op", op),
		slog.Int64("channel_id", s.ChannelID),
		slog.String("created_by", s.CreatedBy),
	)

	// Validation
	err := chh.validator.Struct(s)
	if err != nil {
		log.Warn("invalid parameters", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	msgBody, err := json.Marshal(s)
	if err != nil {
		chh.log.Error("err to marshal shared msg", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = chh.rabbitMQQueues.Publish(ctx, exchangeChannel, channelRoutingKey, msgBody); err != nil {
		chh.log.Error("err send sharing channel to exchange", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("channel sent to share with learning groups")

	return nil
}
