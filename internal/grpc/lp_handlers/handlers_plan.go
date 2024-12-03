package lp_handlers

import (
	"context"
	"errors"
	"time"

	planserv "github.com/DimTur/lp_learning_platform/internal/services/plan"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/plans"
	lpv1 "github.com/DimTur/lp_protos/gen/go/lp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) CreatePlan(ctx context.Context, req *lpv1.CreatePlanRequest) (*lpv1.CreatePlanResponse, error) {
	planID, err := s.planHandlers.CreatePlan(ctx, &plans.CreatePlan{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		CreatedBy:   req.GetCreatedBy(),
		ChannelID:   req.GetChannelId(),
	})
	if err != nil {
		if errors.Is(err, planserv.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lpv1.CreatePlanResponse{
		Id: planID,
	}, nil
}

func (s *serverAPI) GetPlan(ctx context.Context, req *lpv1.GetPlanRequest) (*lpv1.GetPlanResponse, error) {
	plan, err := s.planHandlers.GetPlan(ctx, &plans.GetPlan{
		PlanID:    req.GetPlanId(),
		ChannelID: req.GetChannelId(),
	})
	if err != nil {
		if errors.Is(err, planserv.ErrPlanNotFound) {
			return &lpv1.GetPlanResponse{
				Plan: &lpv1.Plan{},
			}, status.Error(codes.NotFound, "plan not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lpv1.GetPlanResponse{
		Plan: &lpv1.Plan{
			Id:             plan.ID,
			Name:           plan.Name,
			Description:    plan.Description,
			CreatedBy:      plan.CreatedBy,
			LastModifiedBy: plan.LastModifiedBy,
			IsPublished:    plan.IsPublished,
			Public:         plan.Public,
			CreatedAt:      plan.CreatedAt.Format(time.RFC3339),
			Modified:       plan.Modified.Format(time.RFC3339),
		},
	}, nil
}

func (s *serverAPI) GetPlans(ctx context.Context, req *lpv1.GetPlansRequest) (*lpv1.GetPlansResponse, error) {
	plans, err := s.planHandlers.GetPlans(ctx, &plans.GetPlans{
		UserID:    req.GetUserId(),
		ChannelID: req.GetChannelId(),
		Limit:     req.GetLimit(),
		Offset:    req.GetOffset(),
	})
	if err != nil {
		switch {
		case errors.Is(err, planserv.ErrPlanNotFound):
			return nil, status.Error(codes.NotFound, "plans not found")
		case errors.Is(err, planserv.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	var responsePlans []*lpv1.Plan
	for _, plan := range plans {
		responsePlans = append(responsePlans, &lpv1.Plan{
			Id:             plan.ID,
			Name:           plan.Name,
			Description:    plan.Description,
			CreatedBy:      plan.CreatedBy,
			LastModifiedBy: plan.LastModifiedBy,
			IsPublished:    plan.IsPublished,
			Public:         plan.Public,
			CreatedAt:      plan.CreatedAt.Format(time.RFC3339),
			Modified:       plan.Modified.Format(time.RFC3339),
		})
	}

	return &lpv1.GetPlansResponse{
		Plans: responsePlans,
	}, nil
}

func (s *serverAPI) GetPlansAll(ctx context.Context, req *lpv1.GetPlansRequest) (*lpv1.GetPlansResponse, error) {
	plans, err := s.planHandlers.GetPlansAll(ctx, &plans.GetPlans{
		UserID:    req.GetUserId(),
		ChannelID: req.GetChannelId(),
		Limit:     req.GetLimit(),
		Offset:    req.GetOffset(),
	})
	if err != nil {
		switch {
		case errors.Is(err, planserv.ErrPlanNotFound):
			return nil, status.Error(codes.NotFound, "plans not found")
		case errors.Is(err, planserv.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	var responsePlans []*lpv1.Plan
	for _, plan := range plans {
		responsePlans = append(responsePlans, &lpv1.Plan{
			Id:             plan.ID,
			Name:           plan.Name,
			Description:    plan.Description,
			CreatedBy:      plan.CreatedBy,
			LastModifiedBy: plan.LastModifiedBy,
			IsPublished:    plan.IsPublished,
			Public:         plan.Public,
			CreatedAt:      plan.CreatedAt.Format(time.RFC3339),
			Modified:       plan.Modified.Format(time.RFC3339),
		})
	}

	return &lpv1.GetPlansResponse{
		Plans: responsePlans,
	}, nil
}

func (s *serverAPI) UpdatePlan(ctx context.Context, req *lpv1.UpdatePlanRequest) (*lpv1.UpdatePlanResponse, error) {
	id, err := s.planHandlers.UpdatePlan(ctx, &plans.UpdatePlanRequest{
		ChannelID:      req.GetChannelId(),
		PlanID:         req.GetPlanId(),
		Name:           req.GetName(),
		Description:    req.GetDescription(),
		LastModifiedBy: req.GetLastModifiedBy(),
		IsPublished:    req.GetIsPublished(),
		Public:         req.GetPublic(),
	})
	if err != nil {
		switch {
		case errors.Is(err, planserv.ErrPlanNotFound):
			return &lpv1.UpdatePlanResponse{
				Id: 0,
			}, status.Error(codes.NotFound, "plan not found")
		case errors.Is(err, planserv.ErrInvalidCredentials):
			return &lpv1.UpdatePlanResponse{
				Id: 0,
			}, status.Error(codes.InvalidArgument, "bad request")
		default:
			return &lpv1.UpdatePlanResponse{
				Id: 0,
			}, status.Error(codes.Internal, err.Error())
		}
	}

	return &lpv1.UpdatePlanResponse{
		Id: id,
	}, nil
}

func (s *serverAPI) DeletePlan(ctx context.Context, req *lpv1.DeletePlanRequest) (*lpv1.DeletePlanResponse, error) {
	err := s.planHandlers.DeletePlan(ctx, &plans.DeletePlan{
		ChannelID: req.GetChannelId(),
		PlanID:    req.GetPlanId(),
	})
	if err != nil {
		if errors.Is(err, planserv.ErrPlanNotFound) {
			return nil, status.Error(codes.NotFound, "plan not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lpv1.DeletePlanResponse{
		Success: true,
	}, nil
}

func (s *serverAPI) SharePlanWithUsers(ctx context.Context, req *lpv1.SharePlanWithUsersRequest) (*lpv1.SharePlanWithUsersResponse, error) {
	if err := s.planHandlers.SharePlanWithUser(ctx, &plans.SharePlanForUsers{
		ChannelID: req.GetChannelId(),
		PlanID:    req.GetPlanId(),
		UserIDs:   req.GetUsersIds(),
		CreatedBy: req.GetCreatedBy(),
	}); err != nil {
		if errors.Is(err, planserv.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lpv1.SharePlanWithUsersResponse{
		Success: true,
	}, nil
}

func (s *serverAPI) IsUserShareWithPlan(ctx context.Context, req *lpv1.IsUserShareWithPlanRequest) (*lpv1.IsUserShareWithPlanResponse, error) {
	isShare, err := s.planHandlers.IsUserShareWithPlan(ctx, &plans.IsUserShareWithPlan{
		UserID: req.GetUserId(),
		PlanID: req.GetPlanId(),
	})
	if err != nil {
		if errors.Is(err, planserv.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lpv1.IsUserShareWithPlanResponse{
		IsShare: isShare,
	}, nil
}

func (s *serverAPI) GetPlansForSharing(ctx context.Context, req *lpv1.GetPlansForSharingRequest) (*lpv1.GetPlansForSharingResponse, error) {
	resp, err := s.planHandlers.GetPlansForSharing(ctx, &plans.LearningGroup{
		LgID: req.GetLearningGroupId(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var plansForSharing []*lpv1.PlansForSharing
	for channelID, planIDs := range resp {
		plansForSharing = append(plansForSharing, &lpv1.PlansForSharing{
			ChannelId: channelID,
			PlanIds:   planIDs,
		})
	}

	return &lpv1.GetPlansForSharingResponse{
		PlansForSharing: plansForSharing,
	}, nil
}
