package app

import (
	"log/slog"

	grpcapp "github.com/DimTur/lp_learning_platform/internal/app/grpc"
	"github.com/DimTur/lp_learning_platform/internal/services/attempt"
	"github.com/DimTur/lp_learning_platform/internal/services/channel"
	"github.com/DimTur/lp_learning_platform/internal/services/lesson"
	"github.com/DimTur/lp_learning_platform/internal/services/page"
	"github.com/DimTur/lp_learning_platform/internal/services/plan"
	"github.com/DimTur/lp_learning_platform/internal/services/question"
	attstorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/attempts"
	lessonstorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/lessons"
	pagestorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/pages"
	planstorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/plans"
	questiontorage "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/questions"
	"github.com/go-playground/validator/v10"
)

type ChannelStorage interface {
	channel.ChannelSaver
	channel.ChannelProvider
	channel.ChannelDel
}

type ChannelRabbitMq interface {
	channel.RabbitMQQueues
}

type PlanRabbitMq interface {
	plan.RabbitMQQueues
}

type App struct {
	GRPCSrv *grpcapp.Server
}

func NewApp(
	channelStorage ChannelStorage,
	planStorage *planstorage.PlansPostgresStorage,
	lessonStorage *lessonstorage.LessonsPostgresStorage,
	pageStorage *pagestorage.PagesPostgresStorage,
	questionStorage *questiontorage.QuestionsPostgresStorage,
	attemptStorage *attstorage.AttemptsPostgresStorage,
	channelRabbitMq ChannelRabbitMq,
	planRabbitMq PlanRabbitMq,
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
		channelRabbitMq,
	)

	lpGRPCPlanHandlers := plan.New(
		logger,
		validator,
		planStorage,
		planStorage,
		planStorage,
		planRabbitMq,
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

	lpGRPCAttemptHandlers := attempt.New(
		logger,
		validator,
		attemptStorage,
		attemptStorage,
	)

	grpcServer, err := grpcapp.NewGRPCServer(
		grpcAddr,
		lpGRPCChannelHandlers,
		lpGRPCPlanHandlers,
		lpGRPCLessonHandlers,
		lpGRPCPageHandlers,
		lpGRPCQuestionHandlers,
		lpGRPCAttemptHandlers,
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
