package lp_handlers

import (
	"context"

	"github.com/DimTur/lp_learning_platform/internal/services/redis"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/attempts"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/channels"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/lessons"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/pages"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/plans"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/questions"
	lpv1 "github.com/DimTur/lp_protos/gen/go/lp"
	"google.golang.org/grpc"
)

type ChannelHandlers interface {
	CreateChannel(ctx context.Context, channel *channels.CreateChannel) (int64, error)
	GetChannel(ctx context.Context, chLg *channels.GetChannelByID) (*channels.ChannelWithPlans, error)
	GetChannels(ctx context.Context, inputParams *channels.GetChannels) ([]channels.Channel, error)
	UpdateChannel(ctx context.Context, updChannel *channels.UpdateChannelRequest) (int64, error)
	DeleteChannel(ctx context.Context, delChannel *channels.DeleteChannelRequest) error
	ShareChannelToGroup(ctx context.Context, s channels.ShareChannelToGroup) error
	IsChannelCreator(ctx context.Context, isCC *channels.IsChannelCreator) (bool, error)
	GetLearningGroupsShareWithChannel(ctx context.Context, channelID int64) ([]string, error)
}

type PlanHandlers interface {
	CreatePlan(ctx context.Context, plan *plans.CreatePlan) (int64, error)
	GetPlan(ctx context.Context, planCh *plans.GetPlan) (*plans.Plan, error)
	GetPlans(ctx context.Context, inputParams *plans.GetPlans) ([]plans.Plan, error)
	GetPlansAll(ctx context.Context, inputParams *plans.GetPlans) ([]plans.Plan, error)
	UpdatePlan(ctx context.Context, updPlan *plans.UpdatePlanRequest) (int64, error)
	DeletePlan(ctx context.Context, planCh *plans.DeletePlan) error
	SharePlanWithUser(ctx context.Context, s *plans.SharePlanForUsers) error
	IsUserShareWithPlan(ctx context.Context, userPlan *plans.IsUserShareWithPlan) (bool, error)
	GetPlansForSharing(ctx context.Context, lgPlan *plans.LearningGroup) (map[int64][]int64, error)
}

type LessonHandlers interface {
	CreateLesson(ctx context.Context, lesson *lessons.CreateLesson) (int64, error)
	GetLesson(ctx context.Context, lessonPlan *lessons.GetLesson) (*lessons.Lesson, error)
	GetLessons(ctx context.Context, inputParams *lessons.GetLessons) ([]lessons.Lesson, error)
	UpdateLesson(ctx context.Context, updLesson *lessons.UpdateLessonRequest) (int64, error)
	DeleteLesson(ctx context.Context, lessonP *lessons.DeleteLesson) error
}

type PageHandlers interface {
	CreateImagePage(ctx context.Context, imagePage *pages.CreateImagePage) (int64, error)
	CreatePDFPage(ctx context.Context, pdfPage *pages.CreatePDFPage) (int64, error)
	CreateVideoPage(ctx context.Context, videoPage *pages.CreateVideoPage) (int64, error)
	GetImagePage(ctx context.Context, pageLesson *pages.GetPage) (*pages.ImagePage, error)
	GetVideoPage(ctx context.Context, pageLesson *pages.GetPage) (*pages.VideoPage, error)
	GetPDFPage(ctx context.Context, pageLesson *pages.GetPage) (*pages.PDFPage, error)
	GetPages(ctx context.Context, inputParams *pages.GetPages) ([]pages.BasePage, error)
	UpdateImagePage(ctx context.Context, updPage pages.UpdateImagePage) (int64, error)
	UpdateVideoPage(ctx context.Context, updPage pages.UpdateVideoPage) (int64, error)
	UpdatePDFPage(ctx context.Context, updPage pages.UpdatePDFPage) (int64, error)
	DeletePage(ctx context.Context, pageLesson *pages.DeletePage) error
}

type QuestionHandlers interface {
	CreateQuestionPage(ctx context.Context, questionPage *questions.CreateQuestionPage) (int64, error)
	GetQuestionPageByID(ctx context.Context, questionLesson *pages.GetPage) (*questions.QuestionPage, error)
	UpdateQuestionPage(ctx context.Context, updPage *questions.UpdateQuestionPage) (int64, error)
}

type AttemptHandlers interface {
	TryLesson(ctx context.Context, questionPage *attempts.GetQuestionPageAttempts) ([]attempts.QuestionPageAttempt, error)
	UpdatePageAttempt(ctx context.Context, updPAttempt *redis.UpdatePageAttempt) error
	CompleteLesson(ctx context.Context, req *attempts.CompleteLessonRequest) (*attempts.CompleteLessonResp, error)
	GetLessonAttempts(ctx context.Context, inputParams *attempts.GetLessonAttempts) (*attempts.GetLessonAttemptsResp, error)
	CheckPermissionForUser(ctx context.Context, userAtt *attempts.PermissionForUser) (bool, error)
}

type serverAPI struct {
	channelHandlers  ChannelHandlers
	planHandlers     PlanHandlers
	lessonHandlers   LessonHandlers
	pageHandlers     PageHandlers
	questionHandlers QuestionHandlers
	attemptHandlers  AttemptHandlers

	lpv1.UnsafeLearningPlatformServer
}

func RegisterLPServiceServer(
	gRPC *grpc.Server,
	ch ChannelHandlers,
	ph PlanHandlers,
	lh LessonHandlers,
	pgh PageHandlers,
	qh QuestionHandlers,
	ah AttemptHandlers,
) {
	lpv1.RegisterLearningPlatformServer(gRPC, &serverAPI{
		channelHandlers:  ch,
		planHandlers:     ph,
		lessonHandlers:   lh,
		pageHandlers:     pgh,
		questionHandlers: qh,
		attemptHandlers:  ah,
	})
}
