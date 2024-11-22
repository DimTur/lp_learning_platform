package lp_handlers

import (
	"context"
	"errors"
	"time"

	questionserv "github.com/DimTur/lp_learning_platform/internal/services/question"
	pagestore "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/pages"
	questionstore "github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/questions"
	"github.com/DimTur/lp_learning_platform/internal/utils"
	lpv1 "github.com/DimTur/lp_protos/gen/go/lp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) CreateQuestionPage(ctx context.Context, req *lpv1.CreateQuestionPageRequest) (*lpv1.CreateQuestionPageResponse, error) {
	if err := utils.ValidateCreateOptions(req); err != nil {
		return nil, err
	}

	pageID, err := s.questionHandlers.CreateQuestionPage(ctx, &questionstore.CreateQuestionPage{
		LessonID:       req.LessonId,
		CreatedBy:      req.CreatedBy,
		LastModifiedBy: req.CreatedBy,
		ContentType:    "question",
		QuestionType:   "multichoice",
		Question:       req.Question,
		OptionA:        req.OptionA,
		OptionB:        req.OptionB,
		OptionC:        req.GetOptionC(),
		OptionD:        req.GetOptionD(),
		OptionE:        req.GetOptionE(),
		Answer:         req.Answer.String(),
	})
	if err != nil {
		if errors.Is(err, questionserv.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lpv1.CreateQuestionPageResponse{
		Id: pageID,
	}, nil
}

func (s *serverAPI) GetQuestionPage(ctx context.Context, req *lpv1.GetQuestionPageRequest) (*lpv1.GetQuestionPageResponse, error) {
	page, err := s.questionHandlers.GetQuestionPageByID(ctx, &pagestore.GetPage{
		PageID:   req.GetPageId(),
		LessonID: req.GetLessonId(),
	})
	if err != nil {
		if errors.Is(err, questionserv.ErrQuestionNotFound) {
			return nil, status.Error(codes.NotFound, "question page not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lpv1.GetQuestionPageResponse{
		QuestionPage: &lpv1.QuestionPage{
			Id:             page.ID,
			LessonId:       page.LessonID,
			CreatedBy:      page.CreatedBy,
			LastModifiedBy: page.LastModifiedBy,
			CreatedAt:      page.CreatedAt.Format(time.RFC3339),
			Modified:       page.Modified.Format(time.RFC3339),
			ContentType:    lpv1.ContentType_QUESTION,
			QuestionType:   lpv1.QuestionType_MULTICHOICE,
			Question:       page.Question,
			OptionA:        page.OptionA,
			OptionB:        page.OptionB,
			OptionC:        page.OptionC,
			OptionD:        page.OptionD,
			OptionE:        page.OptionE,
			Answer:         page.Answer,
		},
	}, nil
}

func (s *serverAPI) UpdateQuestionPage(ctx context.Context, req *lpv1.UpdateQuestionPageRequest) (*lpv1.UpdateQuestionPageResponse, error) {
	if err := utils.ValidateUpdateOptions(req); err != nil {
		return nil, err
	}

	answer := req.GetAnswer().String()

	id, err := s.questionHandlers.UpdateQuestionPage(ctx, &questionstore.UpdateQuestionPage{
		ID:             req.GetId(),
		LastModifiedBy: req.GetLastModifiedBy(),
		Question:       req.Question,
		OptionA:        req.OptionA,
		OptionB:        req.OptionB,
		OptionC:        req.OptionC,
		OptionD:        req.OptionD,
		OptionE:        req.OptionE,
		Answer:         &answer,
	})
	if err != nil {
		switch {
		case errors.Is(err, questionserv.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		case errors.Is(err, questionserv.ErrQuestionNotFound):
			return nil, status.Error(codes.NotFound, "question page not found")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &lpv1.UpdateQuestionPageResponse{
		Id: id,
	}, nil
}
