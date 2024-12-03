package ssogrpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	ssomodels "github.com/DimTur/lp_learning_platform/internal/clients/sso/models.go"
	ssov1 "github.com/DimTur/lp_protos/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidGroupID     = errors.New("invalid group id")
	ErrGroupExists        = errors.New("group already exists")
	ErrGroupNotFound      = errors.New("group not found")
	ErrPermissionDenied   = errors.New("you don't have permissions")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInternal           = errors.New("internal server error")
	ErrUserNotFound       = errors.New("user not found")
)

func (c *Client) GetLearningGroups(ctx context.Context, uID *ssomodels.GetLGroups) (*ssomodels.GetLGroupsResp, error) {
	const op = "sso.grpc_lg.DeleteLearningGroup"

	lGroups, err := c.api.GetLearningGroups(ctx, &ssov1.GetLearningGroupsRequest{
		UserId: uID.UserID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("groups not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrGroupNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	resp := &ssomodels.GetLGroupsResp{
		LearningGroups: make([]*ssomodels.LearningGroup, len(lGroups.LearningGroups)),
	}
	for i, g := range lGroups.LearningGroups {
		resp.LearningGroups[i] = &ssomodels.LearningGroup{
			Id:         g.Id,
			Name:       g.Name,
			CreatedBy:  g.CreatedBy,
			ModifiedBy: g.ModifiedBy,
			Created:    g.Created,
			Updated:    g.Updated,
		}
	}

	return resp, nil
}

func (c *Client) UserIsGroupAdminIn(ctx context.Context, user *ssomodels.UserIsGroupAdminIn) ([]string, error) {
	const op = "sso.grpc_lg.UserIsGroupAdminIn"

	resp, err := c.api.IsUserGroupAdminIn(ctx, &ssov1.IsUserGroupAdminInRequest{
		UserId: user.UserID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("group not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	groupIDs := append([]string{}, resp.LearningGroupIds...)

	return groupIDs, nil
}

func (c *Client) UserIsLearnerIn(ctx context.Context, user *ssomodels.UserIsLearnerIn) ([]string, error) {
	const op = "sso.grpc_lg.UserIsLearnerIn"

	resp, err := c.api.IsUserLearnerIn(ctx, &ssov1.IsUserLearnereInRequest{
		UserId: user.UserID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("group not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	groupIDs := append([]string{}, resp.LearningGroupIds...)

	return groupIDs, nil
}

func (c *Client) GetLearners(ctx context.Context, lgID string) (*ssomodels.GetLearners, error) {
	const op = "sso.grpc_lg.GetLearnersFromLearningGroup"

	resp, err := c.api.GetLearners(ctx, &ssov1.GetLearnersRequest{
		LearningGroupId: lgID,
	})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.log.Error("learners not found", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		default:
			c.log.Error("internal error", slog.String("err", err.Error()))
			return nil, fmt.Errorf("%s: %w", op, ErrInternal)
		}
	}

	return &ssomodels.GetLearners{
		Learners: resp.Learners,
	}, nil
}
