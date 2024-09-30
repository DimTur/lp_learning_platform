package channel

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/DimTur/lp_learning_platform/internal/domain/models"
	"github.com/DimTur/lp_learning_platform/internal/services/storage"
)

type ChannelSaver interface {
	CreateChannel(ctx context.Context, channel models.Channel) (id int64, err error)
	UpdateChannel(ctx context.Context, updChannel models.UpdateChannelRequest) (id int64, err error)
}

type ChannelProvider interface {
	GetChannelByID(ctx context.Context, channelID int64) (channel models.Channel, err error)
	GetChannels(ctx context.Context, limit, offset int64) (channels []models.Channel, err error)
}

type ChannelDel interface {
	DeleteChannel(ctx context.Context, channelID int64) (err error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidChannelID   = errors.New("invalid channel id")
	ErrChannelExitsts     = errors.New("channel already exists")
	ErrChannelNotFound    = errors.New("channel not found")
)

type LPHandlers struct {
	log             *slog.Logger
	channelSaver    ChannelSaver
	channelProvider ChannelProvider
	channelDel      ChannelDel
}

func New(
	log *slog.Logger,
	channelSaver ChannelSaver,
	channelProvider ChannelProvider,
	channelDel ChannelDel,
) *LPHandlers {
	return &LPHandlers{
		log:             log,
		channelSaver:    channelSaver,
		channelProvider: channelProvider,
		channelDel:      channelDel,
	}
}

// CreateChannel creats new channel in the system and returns channel ID.
func (lp *LPHandlers) CreateChannel(ctx context.Context, channel models.Channel) (int64, error) {
	const op = "channel.CreateChannel"

	log := lp.log.With(
		slog.String("op", op),
		slog.String("name", channel.Name),
	)

	now := time.Now()
	channel.CreatedAt = now
	channel.Modified = now

	log.Info("creating channel")

	id, err := lp.channelSaver.CreateChannel(ctx, channel)
	if err != nil {
		if errors.Is(err, storage.ErrInvalidCredentials) {
			lp.log.Warn("invalid arguments", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to save channel", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// GetChannelByID gets channel by ID and returns it.
func (lp *LPHandlers) GetChannel(ctx context.Context, channelID int64) (models.Channel, error) {
	const op = "channel.GetChannelByID"

	log := lp.log.With(
		slog.String("op", op),
		slog.Int64("chanID", channelID),
	)

	log.Info("getting channel")

	var channel models.Channel
	channel, err := lp.channelProvider.GetChannelByID(ctx, channelID)
	if err != nil {
		if errors.Is(err, storage.ErrChannelNotFound) {
			lp.log.Warn("channel not found", slog.String("err", err.Error()))
			return channel, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to get channel", slog.String("err", err.Error()))
		return channel, fmt.Errorf("%s: %w", op, err)
	}

	return channel, nil
}

// GetChannels gets channels and returns them.
func (lp *LPHandlers) GetChannels(ctx context.Context, limit, offset int64) ([]models.Channel, error) {
	const op = "channel.GetChannels"

	log := lp.log.With(
		slog.String("op", op),
	)

	log.Info("getting channels")

	var channels []models.Channel
	channels, err := lp.channelProvider.GetChannels(ctx, limit, offset)
	if err != nil {
		if errors.Is(err, storage.ErrChannelNotFound) {
			lp.log.Warn("channels not found", slog.String("err", err.Error()))
			return channels, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to get channels", slog.String("err", err.Error()))
		return channels, fmt.Errorf("%s: %w", op, err)
	}

	return channels, nil
}

// UpdateChannel performs a partial update
func (lp *LPHandlers) UpdateChannel(ctx context.Context, updChannel models.UpdateChannelRequest) (int64, error) {
	const op = "channel.UpdateChannel"

	log := lp.log.With(
		slog.String("op", op),
	)

	log.Info("updating channel")

	id, err := lp.channelSaver.UpdateChannel(ctx, updChannel)
	if err != nil {
		if errors.Is(err, storage.ErrChannelNotFound) {
			lp.log.Warn("channel not found", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to update channel", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("channel updated with: ", slog.Int64("channelID", id))

	return id, nil
}

// DeleteChannel
func (lp *LPHandlers) DeleteChannel(ctx context.Context, channelID int64) error {
	const op = "channel.DeleteChannel"

	log := lp.log.With(
		slog.String("op", op),
		slog.Int64("channel id", channelID),
	)

	log.Info("deleting channel with: ", slog.Int64("channelID", channelID))

	err := lp.channelDel.DeleteChannel(ctx, channelID)
	if err != nil {
		if errors.Is(err, storage.ErrChannelNotFound) {
			lp.log.Warn("channel not found", slog.String("err", err.Error()))
			return fmt.Errorf("%s: %w", op, err)
		}

		log.Error("failed to delete channel", slog.String("err", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
