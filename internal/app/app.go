package app

import (
	"log/slog"

	grpcapp "github.com/DimTur/lp_learning_platform/internal/app/grpc"
	"github.com/DimTur/lp_learning_platform/internal/services/channel"
	"github.com/DimTur/lp_learning_platform/internal/services/lesson"
	"github.com/DimTur/lp_learning_platform/internal/services/page"
	"github.com/DimTur/lp_learning_platform/internal/services/plan"
	"github.com/DimTur/lp_learning_platform/internal/services/question"
	channelstorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/channels"
	lessonstorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/lessons"
	pagestorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/pages"
	planstorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/plans"
	questiontorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/questions"
	"github.com/go-playground/validator/v10"
)

type App struct {
	GRPCSrv *grpcapp.Server
}

func NewApp(
	channelStorage *channelstorage.ChannelPostgresStorage,
	planStorage *planstorage.PlansPostgresStorage,
	lessonStorage *lessonstorage.LessonsPostgresStorage,
	pageStorage *pagestorage.PagesPostgresStorage,
	questionStorage *questiontorage.QuestionsPostgresStorage,
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

	lpGRPCPageHandlers := page.New(
		logger,
		validator,
		pageStorage,
		pageStorage,
		pageStorage,
	)

	lpGRPCQuestionHandlers := question.New(
		logger,
		validator,
		questionStorage,
		questionStorage,
	)

	grpcServer, err := grpcapp.NewGRPCServer(
		grpcAddr,
		lpGRPCChannelHandlers,
		lpGRPCPlanHandlers,
		lpGRPCLessonHandlers,
		lpGRPCPageHandlers,
		lpGRPCQuestionHandlers,
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
