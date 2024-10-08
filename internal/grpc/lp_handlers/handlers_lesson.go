package lp_handlers

import (
	"context"
	"errors"

	lessonserv "github.com/DimTur/lp_learning_platform/internal/services/lesson"
	planserv "github.com/DimTur/lp_learning_platform/internal/services/plan"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/lessons"
	lpv1 "github.com/DimTur/lp_learning_platform/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *serverAPI) CreateLesson(ctx context.Context, req *lpv1.CreateLessonRequest) (*lpv1.CreateLessonResponse, error) {
	lesson := lessons.CreateLesson{
		Name:           req.GetName(),
		CreatedBy:      req.GetCreatedBy(),
		LastModifiedBy: req.GetCreatedBy(),
		PlanID:         req.GetPlanId(),
	}

	lessonID, err := s.lessonHandlers.CreateLesson(ctx, lesson)
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
	lesson, err := s.lessonHandlers.GetLesson(ctx, req.GetId())
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
			CreatedAt:      timestamppb.New(lesson.CreatedAt),
			Modified:       timestamppb.New(lesson.Modified),
		},
	}, nil
}

func (s *serverAPI) GetLessons(ctx context.Context, req *lpv1.GetLessonsRequest) (*lpv1.GetLessonsResponse, error) {
	lessons, err := s.lessonHandlers.GetLessons(ctx, req.GetPlanId(), req.GetLimit(), req.GetOffset())
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
			CreatedAt:      timestamppb.New(lesson.CreatedAt),
			Modified:       timestamppb.New(lesson.Modified),
		})
	}

	return &lpv1.GetLessonsResponse{
		Lessons: responseLesson,
	}, nil
}

func (s *serverAPI) UpdateLesson(ctx context.Context, req *lpv1.UpdateLessonRequest) (*lpv1.UpdateLessonResponse, error) {
	var name *string
	if req.GetName() != "" {
		name = proto.String(req.GetName())
	}

	updLesson := lessons.UpdateLessonRequest{
		ID:             req.GetId(),
		Name:           name,
		LastModifiedBy: req.GetLastModifiedBy(),
	}

	id, err := s.lessonHandlers.UpdateLesson(ctx, updLesson)
	if err != nil {
		switch {
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
	lessonID := req.GetId()

	err := s.lessonHandlers.DeleteLesson(ctx, lessonID)
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
