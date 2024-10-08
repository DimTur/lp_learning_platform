package lp_handlers

import (
	"context"
	"errors"

	attserv "github.com/DimTur/lp_learning_platform/internal/services/attempt"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/attempts"
	lpv1 "github.com/DimTur/lp_learning_platform/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) CreateAttempt(ctx context.Context, req *lpv1.CreateAttemptRequest) (*lpv1.CreateAttemptResponse, error) {
	lAttempt := attempts.CreateLessonAttempt{
		LessonID:  req.GetLessonId(),
		PlanId:    req.GetPlanId(),
		ChannelID: req.GetChannelId(),
		UserID:    req.GetUserId(),
	}

	lAttemptID, err := s.attemptHandlers.CreateAttempt(ctx, lAttempt)
	if err != nil {
		if errors.Is(err, attserv.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lpv1.CreateAttemptResponse{
		Id:      lAttemptID,
		Success: true,
	}, nil
}
