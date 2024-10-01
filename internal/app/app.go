package app

import (
	"log/slog"

	grpcapp "github.com/DimTur/lp_learning_platform/internal/app/grpc"
	"github.com/DimTur/lp_learning_platform/internal/services/channel"
	"github.com/DimTur/lp_learning_platform/internal/services/lesson"
	"github.com/DimTur/lp_learning_platform/internal/services/plan"
	channelstorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/channels"
	lessonstorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/lessons"
	planstorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/plans"
	"github.com/go-playground/validator/v10"
)

type App struct {
	GRPCSrv *grpcapp.Server
}

func NewApp(
	channelStorage *channelstorage.ChannelPostgresStorage,
	planStorage *planstorage.PlansPostgresStorage,
	lessonStorage *lessonstorage.LessonsPostgresStorage,
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

	lpGRPCLessonHandlers := lesson.New(
		logger,
		validator,
		lessonStorage,
		lessonStorage,
		lessonStorage,
	)

	grpcServer, err := grpcapp.NewGRPCServer(
		grpcAddr,
		lpGRPCChannelHandlers,
		lpGRPCPlanHandlers,
		lpGRPCLessonHandlers,
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
