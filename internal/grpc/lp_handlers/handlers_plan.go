package lp_handlers

import (
	"context"
	"errors"

	planserv "github.com/DimTur/lp_learning_platform/internal/services/plan"
	"github.com/DimTur/lp_learning_platform/internal/services/storage/postgresql/plans"
	lpv1 "github.com/DimTur/lp_learning_platform/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *serverAPI) CreatePlan(ctx context.Context, req *lpv1.CreatePlanRequest) (*lpv1.CreatePlanResponse, error) {
	plan := plans.CreatePlan{
		Name:           req.GetName(),
		Description:    req.GetDescription(),
		CreatedBy:      req.GetCreatedBy(),
		LastModifiedBy: req.GetCreatedBy(),
		ChannelID:      req.GetChannelId(),
	}

	planID, err := s.planHandlers.CreatePlan(ctx, plan)
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
	plan, err := s.planHandlers.GetPlan(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, planserv.ErrPlanNotFound) {
			return nil, status.Error(codes.NotFound, "plan not found")
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
			CreatedAt:      timestamppb.New(plan.CreatedAt),
			Modified:       timestamppb.New(plan.Modified),
		},
	}, nil
}

func (s *serverAPI) GetPlans(ctx context.Context, req *lpv1.GetPlansRequest) (*lpv1.GetPlansResponse, error) {
	plans, err := s.planHandlers.GetPlans(ctx, req.GetChannelId(), req.GetLimit(), req.GetOffset())
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
			CreatedAt:      timestamppb.New(plan.CreatedAt),
			Modified:       timestamppb.New(plan.Modified),
		})
	}

	return &lpv1.GetPlansResponse{
		Plans: responsePlans,
	}, nil
}

func (s *serverAPI) UpdatePlan(ctx context.Context, req *lpv1.UpdatePlanRequest) (*lpv1.UpdatePlanResponse, error) {
	var name *string
	if req.GetName() != "" {
		name = proto.String(req.GetName())
	}

	var description *string
	if req.GetDescription() != "" {
		description = proto.String(req.GetDescription())
	}

	var isPublished *bool
	if req.IsPublished != nil {
		isPublished = proto.Bool(req.GetIsPublished())
	}

	var public *bool
	if req.Public != nil {
		public = proto.Bool(req.GetPublic())
	}

	updPlan := plans.UpdatePlanRequest{
		ID:             req.GetId(),
		Name:           name,
		Description:    description,
		LastModifiedBy: req.GetLastModifiedBy(),
		IsPublished:    isPublished,
		Public:         public,
	}

	id, err := s.planHandlers.UpdatePlan(ctx, updPlan)
	if err != nil {
		switch {
		case errors.Is(err, planserv.ErrInvalidCredentials):
			return nil, status.Error(codes.InvalidArgument, "bad request")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &lpv1.UpdatePlanResponse{
		Id: id,
	}, nil
}

func (s *serverAPI) DeletePlan(ctx context.Context, req *lpv1.DeletePlanRequest) (*lpv1.DeletePlanResponse, error) {
	planID := req.GetId()

	err := s.planHandlers.DeletePlan(ctx, planID)
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
	sharingPlan := plans.SharePlanForUsers{
		PlanID:    req.GetPlanId(),
		UsersIDs:  req.GetUsersIds(),
		CreatedBy: req.GetCreatedBy(),
	}
	if err := s.planHandlers.SharePlanWithUser(ctx, sharingPlan); err != nil {
		if errors.Is(err, planserv.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &lpv1.SharePlanWithUsersResponse{
		Success: true,
	}, nil
}
