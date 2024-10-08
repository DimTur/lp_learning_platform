package lp_handlers

import (
	"context"

	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/channels"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/lessons"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/pages"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/plans"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/questions"
	lpv1 "github.com/DimTur/lp_learning_platform/pkg/server/grpc"
	"google.golang.org/grpc"
)

type ChannelHandlers interface {
	CreateChannel(ctx context.Context, channel channels.CreateChannel) (int64, error)
	GetChannel(ctx context.Context, channelID int64) (channels.ChannelWithPlans, error)
	GetChannels(ctx context.Context, limit, offset int64) ([]channels.Channel, error)
	UpdateChannel(ctx context.Context, updChannel channels.UpdateChannelRequest) (int64, error)
	DeleteChannel(ctx context.Context, channelID int64) error
}

type PlanHandlers interface {
	CreatePlan(ctx context.Context, plan plans.CreatePlan) (int64, error)
	GetPlan(ctx context.Context, planID int64) (plan plans.Plan, err error)
	GetPlans(ctx context.Context, channel_id int64, limit, offset int64) ([]plans.Plan, error)
	UpdatePlan(ctx context.Context, updPlan plans.UpdatePlanRequest) (int64, error)
	DeletePlan(ctx context.Context, planID int64) error
}

type LessonHandlers interface {
	CreateLesson(ctx context.Context, lesson lessons.CreateLesson) (int64, error)
	GetLesson(ctx context.Context, lessonID int64) (lessons.Lesson, error)
	GetLessons(ctx context.Context, plan_id int64, limit, offset int64) ([]lessons.Lesson, error)
	UpdateLesson(ctx context.Context, updLEsson lessons.UpdateLessonRequest) (int64, error)
	DeleteLesson(ctx context.Context, lessonID int64) error
}

type PageHandlers interface {
	CreatePage(ctx context.Context, page pages.CreatePage) (int64, error)
	GetPage(ctx context.Context, pageID int64, contentType string) (pages.Page, error)
	GetPages(ctx context.Context, lessonID int64, limit, offset int64) ([]pages.BasePage, error)
	UpdatePage(ctx context.Context, updPage pages.UpdatePage) (int64, error)
	DeletePage(ctx context.Context, pageID int64) error
}

type QuestionHandlers interface {
	CreateQuestionPage(ctx context.Context, questionPage questions.CreateQuestionPage) (int64, error)
	GetQuestionPageByID(ctx context.Context, pageID int64) (questions.QuestionPage, error)
	UpdateQuestionPage(ctx context.Context, updPage questions.UpdateQuestionPage) (int64, error)
}

type serverAPI struct {
	channelHandlers  ChannelHandlers
	planHandlers     PlanHandlers
	lessonHandlers   LessonHandlers
	pageHandlers     PageHandlers
	questionHandlers QuestionHandlers

	lpv1.UnsafeLearningPlatformServer
}

func RegisterLPServiceServer(
	gRPC *grpc.Server,
	ch ChannelHandlers,
	ph PlanHandlers,
	lh LessonHandlers,
	pgh PageHandlers,
	qh QuestionHandlers,
) {
	lpv1.RegisterLearningPlatformServer(gRPC, &serverAPI{
		channelHandlers:  ch,
		planHandlers:     ph,
		lessonHandlers:   lh,
		pageHandlers:     pgh,
		questionHandlers: qh,
	})
}
