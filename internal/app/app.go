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
	"github.com/go-playground/validator/v10"
)

type ChannelStorage interface {
	channel.ChannelSaver
	channel.ChannelProvider
	channel.ChannelDel
}

type PlanlStorage interface {
	plan.PlanSaver
	plan.PlanProvider
	plan.PlanDel
}

type LessonStorage interface {
	lesson.LessonSaver
	lesson.LessonProvider
	lesson.LessonDel
}

type PageStorage interface {
	page.PageSaver
	page.PageProvider
	page.PageDel
}

type QuestionStorage interface {
	question.QuestionPageSaver
	question.QuestionPageProvider
}

type AttemptStorage interface {
	attempt.AttemptSaver
	attempt.AttemptProvider
}

type ChannelRabbitMq interface {
	channel.RabbitMQQueues
}

type PlanRabbitMq interface {
	plan.RabbitMQQueues
}

type AttemptsRedis interface {
	attempt.AttemptRedisStore
}

type SsoStorage interface {
	plan.LearningGroupProvider
}

type App struct {
	GRPCSrv *grpcapp.Server
}

func NewApp(
	channelStorage ChannelStorage,
	planStorage PlanlStorage,
	lessonStorage LessonStorage,
	pageStorage PageStorage,
	questionStorage QuestionStorage,
	attemptStorage AttemptStorage,
	attemptRedis AttemptsRedis,
	channelRabbitMq ChannelRabbitMq,
	planRabbitMq PlanRabbitMq,
	ssoStorage SsoStorage,
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
		channelStorage,
		ssoStorage,
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
		attemptRedis,
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
