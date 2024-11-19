package lp_handlers

import (
	"context"
	"errors"
	"time"

	lessonserv "github.com/DimTur/lp_learning_platform/internal/services/lesson"
	planserv "github.com/DimTur/lp_learning_platform/internal/services/plan"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/lessons"
	lpv1 "github.com/DimTur/lp_protos/gen/go/lp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) CreateLesson(ctx context.Context, req *lpv1.CreateLessonRequest) (*lpv1.CreateLessonResponse, error) {

	lessonID, err := s.lessonHandlers.CreateLesson(ctx, &lessons.CreateLesson{
		Name:           req.GetName(),
		Description:    req.GetDescription(),
		CreatedBy:      req.GetCreatedBy(),
		LastModifiedBy: req.GetCreatedBy(),
		PlanID:         req.GetPlanId(),
	})
	if err != nil {
		if errors.Is(err, planserv.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lpv1.CreateLessonResponse{
		Id: lessonID,
	}, nil
}

func (s *serverAPI) GetLesson(ctx context.Context, req *lpv1.GetLessonRequest) (*lpv1.GetLessonResponse, error) {
	lesson, err := s.lessonHandlers.GetLesson(ctx, &lessons.GetLesson{
		LessonID: req.LessonId,
		PlanID:   req.PlanId,
	})
	if err != nil {
		if errors.Is(err, lessonserv.ErrLessonNotFound) {
			return nil, status.Error(codes.NotFound, "lesson not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lpv1.GetLessonResponse{
		Lesson: &lpv1.Lesson{
			Id:             lesson.ID,
			Name:           lesson.Name,
			CreatedBy:      lesson.CreatedBy,
			LastModifiedBy: lesson.LastModifiedBy,
			CreatedAt:      lesson.CreatedAt.Format(time.RFC3339),
			Modified:       lesson.Modified.Format(time.RFC3339),
		},
	}, nil
}

func (s *serverAPI) GetLessons(ctx context.Context, req *lpv1.GetLessonsRequest) (*lpv1.GetLessonsResponse, error) {
	lessons, err := s.lessonHandlers.GetLessons(ctx, &lessons.GetLessons{
		PlanID: req.GetPlanId(),
		Limit:  req.GetLimit(),
		Offset: req.GetOffset(),
	})
	if err != nil {
		switch {
		case errors.Is(err, lessonserv.ErrLessonNotFound):
			return nil, status.Error(codes.NotFound, "lessons not found")
		case errors.Is(err, lessonserv.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	var responseLesson []*lpv1.Lesson
	for _, lesson := range lessons {
		responseLesson = append(responseLesson, &lpv1.Lesson{
			Id:             lesson.ID,
			Name:           lesson.Name,
			CreatedBy:      lesson.CreatedBy,
			LastModifiedBy: lesson.LastModifiedBy,
			CreatedAt:      lesson.CreatedAt.Format(time.RFC3339),
			Modified:       lesson.Modified.Format(time.RFC3339),
		})
	}

	return &lpv1.GetLessonsResponse{
		Lessons: responseLesson,
	}, nil
}

func (s *serverAPI) UpdateLesson(ctx context.Context, req *lpv1.UpdateLessonRequest) (*lpv1.UpdateLessonResponse, error) {
	id, err := s.lessonHandlers.UpdateLesson(ctx, &lessons.UpdateLessonRequest{
		PlanID:         req.GetPlanId(),
		LessonID:       req.GetLessonId(),
		Name:           req.GetName(),
		Description:    req.GetDescription(),
		LastModifiedBy: req.GetLastModifiedBy(),
	})
	if err != nil {
		switch {
		case errors.Is(err, lessonserv.ErrLessonNotFound):
			return nil, status.Error(codes.NotFound, "lesson not found")
		case errors.Is(err, lessonserv.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &lpv1.UpdateLessonResponse{
		Id: id,
	}, nil
}

func (s *serverAPI) DeleteLesson(ctx context.Context, req *lpv1.DeleteLessonRequest) (*lpv1.DeleteLessonResponse, error) {
	err := s.lessonHandlers.DeleteLesson(ctx, &lessons.DeleteLesson{
		LessonID: req.LessonId,
		PlanID:   req.PlanId,
	})
	if err != nil {
		if errors.Is(err, lessonserv.ErrLessonNotFound) {
			return nil, status.Error(codes.NotFound, "lesson not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lpv1.DeleteLessonResponse{
		Success: true,
	}, nil
}
