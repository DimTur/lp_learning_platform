package lp_handlers

import (
	"context"
	"errors"

	"github.com/DimTur/lp_learning_platform/internal/domain/models"
	chanserv "github.com/DimTur/lp_learning_platform/internal/services/channel"
	lpv1 "github.com/DimTur/lp_learning_platform/pkg/server/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChannelHandlers interface {
	CreateChannel(ctx context.Context, channel models.Channel) (id int64, err error)
	GetChannel(ctx context.Context, channelID int64) (channel models.Channel, err error)
	GetChannels(ctx context.Context, limit, offset int64) (channels []models.Channel, err error)
	UpdateChannel(ctx context.Context, updChannel models.UpdateChannelRequest) (id int64, err error)
	DeleteChannel(ctx context.Context, channelID int64) (err error)
}

type serverAPI struct {
	channelHandlers ChannelHandlers

	lpv1.UnsafeLearningPlatformServer
}

func RegisterLPServiceServer(gRPC *grpc.Server, ch ChannelHandlers) {
	lpv1.RegisterLearningPlatformServer(gRPC, &serverAPI{channelHandlers: ch})
}

func (s *serverAPI) CreateChannel(ctx context.Context, req *lpv1.CreateChannelRequest) (*lpv1.CreateChannelResponse, error) {
	reqChan := req.GetChannel()

	channel := models.Channel{
		Name:           reqChan.GetName(),
		Description:    reqChan.GetDescription(),
		CreatedBy:      reqChan.GetCreatedBy(),
		LastModifiedBy: reqChan.GetCreatedBy(),
	}

	channelID, err := s.channelHandlers.CreateChannel(ctx, channel)
	if err != nil {
		if errors.Is(err, chanserv.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lpv1.CreateChannelResponse{
		Channel: &lpv1.Channel{
			Id:             channelID,
			Name:           channel.Name,
			Description:    channel.Description,
			CreatedBy:      channel.CreatedBy,
			LastModifiedBy: channel.LastModifiedBy,
		},
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

	return &lpv1.GetChannelResponse{
		Channel: &lpv1.Channel{
			Id:             channel.ID,
			Name:           channel.Name,
			Description:    channel.Description,
			CreatedBy:      channel.CreatedBy,
			LastModifiedBy: channel.LastModifiedBy,
			CreatedAt:      timestamppb.New(channel.CreatedAt),
			Modified:       timestamppb.New(channel.Modified),
		},
	}, nil
}

func (s *serverAPI) GetChannels(ctx context.Context, req *lpv1.GetChannelsRequest) (*lpv1.GetChannelsResponse, error) {
	channels, err := s.channelHandlers.GetChannels(ctx, req.GetLimit(), req.GetOffset())
	if err != nil {
		if errors.Is(err, chanserv.ErrChannelNotFound) {
			return nil, status.Error(codes.NotFound, "channels not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
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
	if req.GetChannel().GetName() != "" {
		name = proto.String(req.GetChannel().GetName())
	}

	var description *string
	if req.GetChannel().GetDescription() != "" {
		description = proto.String(req.GetChannel().GetDescription())
	}

	updChannel := models.UpdateChannelRequest{
		ID:             req.GetChannel().GetId(),
		Name:           name,
		Description:    description,
		LastModifiedBy: req.GetChannel().GetLastModifiedBy(),
	}

	id, err := s.channelHandlers.UpdateChannel(ctx, updChannel)
	if err != nil {
		if errors.Is(err, chanserv.ErrChannelNotFound) {
			return nil, status.Error(codes.NotFound, "channel not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lpv1.UpdateChannelResponse{
		Channel: &lpv1.UpdateChannel{
			Id:             id,
			Name:           req.GetChannel().GetName(),
			Description:    req.GetChannel().GetDescription(),
			LastModifiedBy: req.GetChannel().GetLastModifiedBy(),
		},
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
