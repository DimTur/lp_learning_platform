package lp_handlers

import (
	"context"
	"errors"

	"github.com/DimTur/lp_learning_platform/internal/domain/models"
	chanserv "github.com/DimTur/lp_learning_platform/internal/services/channel"
	lpv1 "github.com/DimTur/lp_protos/gen/go/lp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LPHandlers interface {
	CreateChannel(ctx context.Context,
		name string,
		description string,
		userID int64,
		public bool) (id int64, err error)
	// CreatePlan(ctx context.Context, plan lpv1.CreatePlanRequest) (resp lpv1.CreatePlanResponse, err error)
	// CreateLesson(ctx context.Context, lesson lpv1.CreateLessonRequest) (resp lpv1.CreateLessonResponse, err error)
	GetChannel(ctx context.Context, channelID int64) (channel models.Channel, err error)
	// GetPlan(ctx context.Context, plan lpv1.GetPlanRequest) (resp lpv1.GetPlanResponse, err error)
	// GetLesson(ctx context.Context, lesson lpv1.GetLessonRequest) (resp lpv1.GetLessonResponse, err error)
}

type serverAPI struct {
	learningPlatform LPHandlers

	lpv1.UnsafeLearningPlatformServer
}

func RegisterLPServiceServer(gRPC *grpc.Server, lp LPHandlers) {
	lpv1.RegisterLearningPlatformServer(gRPC, &serverAPI{learningPlatform: lp})
}

func (s *serverAPI) CreateChannel(ctx context.Context, req *lpv1.CreateChannelRequest) (*lpv1.CreateChannelResponse, error) {
	reqChan := req.GetChannel()
	channelID, err := s.learningPlatform.CreateChannel(
		ctx,
		reqChan.GetName(),
		reqChan.GetDescription(),
		reqChan.GetCreatedBy(),
		reqChan.GetPublic(),
	)
	if err != nil {
		if errors.Is(err, chanserv.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.InvalidArgument, "invalid credentials")
	}

	return &lpv1.CreateChannelResponse{
		Channel: &lpv1.Channel{
			ChannelId:      channelID,
			Name:           reqChan.GetName(),
			Description:    reqChan.GetDescription(),
			CreatedBy:      reqChan.GetCreatedBy(),
			Public:         reqChan.GetPublic(),
			LastModifiedBy: reqChan.GetCreatedBy(),
		},
	}, nil
}

func (s *serverAPI) CreatePlan(ctx context.Context, req *lpv1.CreatePlanRequest) (*lpv1.CreatePlanResponse, error) {
	return &lpv1.CreatePlanResponse{}, nil
}

func (s *serverAPI) CreateLesson(ctx context.Context, req *lpv1.CreateLessonRequest) (*lpv1.CreateLessonResponse, error) {
	return &lpv1.CreateLessonResponse{}, nil
}

func (s *serverAPI) GetChannel(ctx context.Context, req *lpv1.GetChannelRequest) (*lpv1.GetChannelResponse, error) {
	channel, err := s.learningPlatform.GetChannel(ctx, req.GetChannelId())
	if err != nil {
		if errors.Is(err, chanserv.ErrInvalidCredentials) {
			return nil, status.Error(codes.NotFound, "channel not found")
		}

		return nil, status.Error(codes.NotFound, "channel not found")
	}

	return &lpv1.GetChannelResponse{
		Channel: &lpv1.Channel{
			ChannelId:      channel.ID,
			Name:           channel.Name,
			Description:    channel.Description,
			CreatedBy:      channel.CreatedBy,
			Public:         channel.Public,
			LastModifiedBy: channel.LastModifiedBy,
		},
	}, nil
}

func (s *serverAPI) GetPlan(ctx context.Context, req *lpv1.GetPlanRequest) (*lpv1.GetPlanResponse, error) {
	return &lpv1.GetPlanResponse{}, nil
}

func (s *serverAPI) GetLesson(ctx context.Context, req *lpv1.GetLessonRequest) (*lpv1.GetLessonResponse, error) {
	return &lpv1.GetLessonResponse{}, nil
}
