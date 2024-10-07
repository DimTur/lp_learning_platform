package lp_handlers

import (
	"context"
	"errors"

	"github.com/DimTur/lp_learning_platform/internal/domain/models"
	chanserv "github.com/DimTur/lp_learning_platform/internal/services/channel"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/pages"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/questions"
	lpv1 "github.com/DimTur/lp_learning_platform/pkg/server/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChannelHandlers interface {
	CreateChannel(ctx context.Context, channel models.CreateChannel) (id int64, err error)
	GetChannel(ctx context.Context, channelID int64) (channel models.ChannelWithPlans, err error)
	GetChannels(ctx context.Context, limit, offset int64) (channels []models.Channel, err error)
	UpdateChannel(ctx context.Context, updChannel models.UpdateChannelRequest) (id int64, err error)
	DeleteChannel(ctx context.Context, channelID int64) (err error)
}

type PlanHandlers interface {
	CreatePlan(ctx context.Context, plan models.CreatePlan) (id int64, err error)
	GetPlan(ctx context.Context, planID int64) (plan models.Plan, err error)
	GetPlans(ctx context.Context, channel_id int64, limit, offset int64) (plans []models.Plan, err error)
	UpdatePlan(ctx context.Context, updPlan models.UpdatePlanRequest) (id int64, err error)
	DeletePlan(ctx context.Context, planID int64) (err error)
}

type LessonHandlers interface {
	CreateLesson(ctx context.Context, lesson models.CreateLesson) (id int64, err error)
	GetLesson(ctx context.Context, lessonID int64) (lesson models.Lesson, err error)
	GetLessons(ctx context.Context, plan_id int64, limit, offset int64) (lessons []models.Lesson, err error)
	UpdateLesson(ctx context.Context, updLEsson models.UpdateLessonRequest) (id int64, err error)
	DeleteLesson(ctx context.Context, lessonID int64) (err error)
}

type PageHandlers interface {
	CreatePage(ctx context.Context, page pages.CreatePage) (int64, error)
	GetPage(ctx context.Context, pageID int64, contentType string) (pages.Page, error)
	GetPages(ctx context.Context, lessonID int64, limit, offset int64) ([]pages.BasePage, error)
	UpdatePage(ctx context.Context, updPage pages.UpdatePage) (int64, error)
	DeletePage(ctx context.Context, pageID int64) error
	// CreateImagePage(ctx context.Context, imagePage models.CreateImagePage) (int64, error)
	// CreateVideoPage(ctx context.Context, videoPage models.CreateVideoPage) (int64, error)
	// CreatePDFPage(ctx context.Context, pdfPage models.CreatePDFPage) (int64, error)
	// CreateQuestionPage(ctx context.Context, questionPage models.CreateQuestionPage) (int64, error)
	// GetImagePageByID(ctx context.Context, pageID int64) (models.ImagePage, error)
	// GetVideoPageByID(ctx context.Context, pageID int64) (models.VideoPage, error)
	// GetPDFPageByID(ctx context.Context, pageID int64) (models.PDFPage, error)
	// GetQuestionPageByID(ctx context.Context, pageID int64) (models.QuestionPage, error)
	// GetPages(ctx context.Context, lessonID int64, limit, offset int64) ([]models.Page, error)
	// UpdateImagePage(ctx context.Context, updPage models.UpdateImagePage) (int64, error)
	// UpdateVideoPage(ctx context.Context, updPage models.UpdateVideoPage) (int64, error)
	// UpdatePDFPage(ctx context.Context, updPage models.UpdatePDFPage) (int64, error)
	// UpdateQuestionPage(ctx context.Context, updPage models.UpdateQuestionPage) (int64, error)
	// DeletePage(ctx context.Context, pageID int64) error
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

func (s *serverAPI) CreateChannel(ctx context.Context, req *lpv1.CreateChannelRequest) (*lpv1.CreateChannelResponse, error) {
	channel := models.CreateChannel{
		Name:           req.GetName(),
		Description:    req.GetDescription(),
		CreatedBy:      req.GetCreatedBy(),
		LastModifiedBy: req.GetCreatedBy(),
	}

	channelID, err := s.channelHandlers.CreateChannel(ctx, channel)
	if err != nil {
		if errors.Is(err, chanserv.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lpv1.CreateChannelResponse{
		Id: channelID,
	}, nil
}

func (s *serverAPI) GetChannel(ctx context.Context, req *lpv1.GetChannelRequest) (*lpv1.GetChannelResponse, error) {
	channel, err := s.channelHandlers.GetChannel(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, chanserv.ErrChannelNotFound) {
			return nil, status.Error(codes.NotFound, "channel not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	var plans []*lpv1.Plan
	for _, plan := range channel.Plans {
		plans = append(plans, &lpv1.Plan{
			Id:             plan.ID,
			Name:           plan.Name,
			Description:    plan.Description,
			CreatedBy:      plan.CreatedBy,
			LastModifiedBy: plan.LastModifiedBy,
			IsPublished:    plan.IsPublished,
			Public:         plan.Public,
			CreatedAt:      timestamppb.New(plan.CreatedAt),
			Modified:       timestamppb.New(plan.Modified),
		})
	}

	return &lpv1.GetChannelResponse{
		Channel: &lpv1.ChannelWithPlans{
			Id:             channel.ID,
			Name:           channel.Name,
			Description:    channel.Description,
			CreatedBy:      channel.CreatedBy,
			LastModifiedBy: channel.LastModifiedBy,
			CreatedAt:      timestamppb.New(channel.CreatedAt),
			Modified:       timestamppb.New(channel.Modified),
			Plans:          plans,
		},
	}, nil
}

func (s *serverAPI) GetChannels(ctx context.Context, req *lpv1.GetChannelsRequest) (*lpv1.GetChannelsResponse, error) {
	channels, err := s.channelHandlers.GetChannels(ctx, req.GetLimit(), req.GetOffset())
	if err != nil {
		switch {
		case errors.Is(err, chanserv.ErrChannelNotFound):
			return nil, status.Error(codes.NotFound, "channels not found")
		case errors.Is(err, chanserv.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	var responseChannels []*lpv1.Channel
	for _, channel := range channels {
		responseChannels = append(responseChannels, &lpv1.Channel{
			Id:             channel.ID,
			Name:           channel.Name,
			Description:    channel.Description,
			CreatedBy:      channel.CreatedBy,
			LastModifiedBy: channel.LastModifiedBy,
			CreatedAt:      timestamppb.New(channel.CreatedAt),
			Modified:       timestamppb.New(channel.Modified),
		})
	}

	return &lpv1.GetChannelsResponse{
		Channels: responseChannels,
	}, nil
}

func (s *serverAPI) UpdateChannel(ctx context.Context, req *lpv1.UpdateChannelRequest) (*lpv1.UpdateChannelResponse, error) {
	var name *string
	if req.GetName() != "" {
		name = proto.String(req.GetName())
	}

	var description *string
	if req.GetDescription() != "" {
		description = proto.String(req.GetDescription())
	}

	updChannel := models.UpdateChannelRequest{
		ID:             req.GetId(),
		Name:           name,
		Description:    description,
		LastModifiedBy: req.GetLastModifiedBy(),
	}

	id, err := s.channelHandlers.UpdateChannel(ctx, updChannel)
	if err != nil {
		switch {
		case errors.Is(err, chanserv.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &lpv1.UpdateChannelResponse{
		Id: id,
	}, nil
}

func (s *serverAPI) DeleteChannel(ctx context.Context, req *lpv1.DeleteChannelRequest) (*lpv1.DeleteChannelResponse, error) {
	channelID := req.GetId()

	err := s.channelHandlers.DeleteChannel(ctx, channelID)
	if err != nil {
		if errors.Is(err, chanserv.ErrChannelNotFound) {
			return nil, status.Error(codes.NotFound, "channel not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lpv1.DeleteChannelResponse{
		Success: true,
	}, nil
}
