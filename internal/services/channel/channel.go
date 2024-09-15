package channel

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/DimTur/lp_learning_platform/internal/domain/models"
	"github.com/DimTur/lp_learning_platform/internal/services/storage"
)

type ChannelSaver interface {
	SaveChannel(
		ctx context.Context,
		name string,
		description string,
		userID int64,
		public bool,
	) (id int64, err error)
}

type ChannelProvider interface {
	GetChannelByID(ctx context.Context, channelID int64) (channel models.Channel, err error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidChannelID   = errors.New("invalid channel id")
	ErrChannelExitsts     = errors.New("channel already exists")
)

type LPHandlers struct {
	log             *slog.Logger
	channelSaver    ChannelSaver
	channelProvider ChannelProvider
}

func New(
	log *slog.Logger,
	channelSaver ChannelSaver,
	channelProvider ChannelProvider,
) *LPHandlers {
	return &LPHandlers{
		log:             log,
		channelSaver:    channelSaver,
		channelProvider: channelProvider,
	}
}

// CreateChannel creats new channel in the system and returns channel ID.
func (lp *LPHandlers) CreateChannel(ctx context.Context, name string, description string, userID int64, public bool) (int64, error) {
	const op = "channel.CreateChannel"

	log := lp.log.With(
		slog.String("op", op),
		slog.String("name", name),
	)

	log.Info("creating channel")

	id, err := lp.channelSaver.SaveChannel(ctx, name, description, userID, public)
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
func (lp *LPHandlers) GetChannel(ctx context.Context, chanID int64) (models.Channel, error) {
	const op = "channel.GetChannelByID"

	log := lp.log.With(
		slog.String("op", op),
		slog.Int64("chanID", chanID),
	)

	log.Info("getting channel")

	var channel models.Channel
	channel, err := lp.channelProvider.GetChannelByID(ctx, chanID)
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
