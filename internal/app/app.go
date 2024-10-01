package app

import (
	"log/slog"

	grpcapp "github.com/DimTur/lp_learning_platform/internal/app/grpc"
	"github.com/DimTur/lp_learning_platform/internal/services/channel"
	"github.com/DimTur/lp_learning_platform/internal/services/plan"
	channelstorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/channels"
	planstorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/plans"
	"github.com/go-playground/validator/v10"
)

type App struct {
	GRPCSrv *grpcapp.Server
}

func NewApp(
	channelStorage *channelstorage.ChannelPostgresStorage,
	planStorage *planstorage.PlansPostgresStorage,
	grpcAddr string,
	logger *slog.Logger,
	validator *validator.Validate,
) (*App, error) {
	lpGRPCChannelHandlers := channel.New(
		logger,
		validator,
		channelStorage,
		channelStorage,
		channelStorage,
	)

	lpGRPCPlanHandlers := plan.New(
		logger,
		validator,
		planStorage,
		planStorage,
		planStorage,
	)

	grpcServer, err := grpcapp.NewGRPCServer(
		grpcAddr,
		lpGRPCChannelHandlers,
		lpGRPCPlanHandlers,
		logger,
		validator,
	)
	if err != nil {
		return nil, err
	}

	return &App{
		GRPCSrv: grpcServer,
	}, nil
}
