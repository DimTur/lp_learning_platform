package lp_handlers

import (
	"context"
	"errors"
	"time"

	attemptserve "github.com/DimTur/lp_learning_platform/internal/services/attempt"
	"github.com/DimTur/lp_learning_platform/internal/services/redis"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/attempts"
	lpv1 "github.com/DimTur/lp_protos/gen/go/lp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) TryLesson(ctx context.Context, req *lpv1.TryLessonRequest) (*lpv1.TryLessonResponse, error) {
	qPAttempts, err := s.attemptHandlers.TryLesson(ctx, &attempts.GetQuestionPageAttempts{
		UserID:    req.GetUserId(),
		LessonID:  req.GetLessonId(),
		PlanID:    req.GetPlanId(),
		ChannelID: req.GetChannelId(),
	})
	if err != nil {
		switch {
		case errors.Is(err, attemptserve.ErrPageAttemtsNotFound):
			return nil, status.Error(codes.NotFound, "question page attempts not found")
		case errors.Is(err, attemptserve.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	var resp []*lpv1.QuestionPageAttempt
	for _, qPAttempt := range qPAttempts {
		resp = append(resp, &lpv1.QuestionPageAttempt{
			Id:              qPAttempt.ID,
			PageId:          qPAttempt.PageID,
			LessonAttemptId: qPAttempt.LessonAttemptID,
			IsCorrect:       qPAttempt.IsCorrect,
			UserAnswer:      stringToAnswer(qPAttempt.UserAnswer),
		})
	}

	return &lpv1.TryLessonResponse{
		QuestionPageAttempts: resp,
	}, nil
}

func (s *serverAPI) UpdatePageAttempt(ctx context.Context, req *lpv1.UpdatePageAttemptRequest) (*lpv1.UpdatePageAttemptResponse, error) {
	if err := s.attemptHandlers.UpdatePageAttempt(ctx, &redis.UpdatePageAttempt{
		LessonAttemptID: req.GetLessonAttemptId(),
		PageID:          req.GetPageId(),
		QPAttemptID:     req.GetQuestionAttemptId(),
		UserAnswer:      req.GetUserAnswer().Enum().String(),
	}); err != nil {
		switch {
		case errors.Is(err, attemptserve.ErrAnswerNotFound):
			return nil, status.Error(codes.NotFound, "answer not found")
		case errors.Is(err, attemptserve.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &lpv1.UpdatePageAttemptResponse{
		Success: true,
	}, nil
}

func (s *serverAPI) CompleteLesson(ctx context.Context, req *lpv1.CompleteLessonRequest) (*lpv1.CompleteLessonResponse, error) {
	resp, err := s.attemptHandlers.CompleteLesson(ctx, &attempts.CompleteLessonRequest{
		UserID:          req.GetUserId(),
		LessonAttemptID: req.GetLessonAttemptId(),
	})
	if err != nil {
		switch {
		case errors.Is(err, attemptserve.ErrLessonAttemtNotFound):
			return nil, status.Error(codes.NotFound, "lesson attempt not found")
		case errors.Is(err, attemptserve.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &lpv1.CompleteLessonResponse{
		LessonAttemptId: resp.ID,
		IsSuccessfull:   resp.IsSuccessful,
		PercentageScore: resp.PercentageScore,
	}, nil
}

func (s *serverAPI) GetLessonAttempts(ctx context.Context, req *lpv1.GetLessonAttemptsRequest) (*lpv1.GetLessonAttemptsResponse, error) {
	resp, err := s.attemptHandlers.GetLessonAttempts(ctx, &attempts.GetLessonAttempts{
		UserID:   req.GetUserId(),
		LessonID: req.GetLessonId(),
		Limit:    req.GetLimit(),
		Offset:   req.GetOffset(),
	})
	if err != nil {
		switch {
		case errors.Is(err, attemptserve.ErrLessonAttemtNotFound):
			return nil, status.Error(codes.NotFound, "lesson attempt not found")
		case errors.Is(err, attemptserve.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	var responseAttempt []*lpv1.LessonAttempt
	for _, attempt := range resp.LessonAttempts {
		responseAttempt = append(responseAttempt, &lpv1.LessonAttempt{
			Id:              attempt.ID,
			UserId:          attempt.UserID,
			LessonId:        attempt.LessonID,
			PlanId:          attempt.PlanID,
			ChannelId:       attempt.ChannelID,
			StartTime:       attempt.StartTime.Format(time.RFC3339),
			EndTime:         attempt.EndTime.Format(time.RFC3339),
			IsComplete:      attempt.IsComplete,
			IsSuccessful:    attempt.IsSuccessful,
			PercentageScore: attempt.PercentageScore,
		})
	}

	return &lpv1.GetLessonAttemptsResponse{
		LessonAttempts: responseAttempt,
	}, nil
}

func (s *serverAPI) CheckPermissionForUser(ctx context.Context, req *lpv1.CheckPermissionForUserRequest) (*lpv1.CheckPermissionForUserResponse, error) {
	resp, err := s.attemptHandlers.CheckPermissionForUser(ctx, &attempts.PermissionForUser{
		UserID:          req.GetUserId(),
		LessonAttemptID: req.GetLessonAttemptId(),
	})
	if err != nil {
		switch {
		case errors.Is(err, attemptserve.ErrPermissionsDenied):
			return nil, status.Error(codes.PermissionDenied, "permissions denied")
		case errors.Is(err, attemptserve.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &lpv1.CheckPermissionForUserResponse{
		Success: resp,
	}, nil
}

func stringToAnswer(answer string) lpv1.Answer {
	switch answer {
	case "OPTION_A":
		return lpv1.Answer_OPTION_A
	case "OPTION_B":
		return lpv1.Answer_OPTION_B
	case "OPTION_C":
		return lpv1.Answer_OPTION_C
	case "OPTION_D":
		return lpv1.Answer_OPTION_D
	case "OPTION_E":
		return lpv1.Answer_OPTION_E
	default:
		return lpv1.Answer_ANSWER_UNSPECIFIED
	}
}
