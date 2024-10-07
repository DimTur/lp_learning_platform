package utils

import (
	lpv1 "github.com/DimTur/lp_learning_platform/pkg/server/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidateCreateOptions(req *lpv1.CreateQuestionPageRequest) error {
	switch req.GetAnswer() {
	case lpv1.Answer_OPTION_A:
		if req.GetOptionA() == "" {
			return status.Error(codes.InvalidArgument, "option A must be provided if it is selected as the answer")
		}
	case lpv1.Answer_OPTION_B:
		if req.GetOptionB() == "" {
			return status.Error(codes.InvalidArgument, "option B must be provided if it is selected as the answer")
		}
	case lpv1.Answer_OPTION_C:
		if req.GetOptionC() == "" {
			return status.Error(codes.InvalidArgument, "option C must be provided if it is selected as the answer")
		}
	case lpv1.Answer_OPTION_D:
		if req.GetOptionD() == "" {
			return status.Error(codes.InvalidArgument, "option D must be provided if it is selected as the answer")
		}
	case lpv1.Answer_OPTION_E:
		if req.GetOptionE() == "" {
			return status.Error(codes.InvalidArgument, "option E must be provided if it is selected as the answer")
		}
	}

	return nil
}

func ValidateUpdateOptions(req *lpv1.UpdateQuestionPageRequest) error {
	switch req.GetAnswer() {
	case lpv1.Answer_OPTION_A:
		if req.GetOptionA() == "" {
			return status.Error(codes.InvalidArgument, "option A must be provided if it is selected as the answer")
		}
	case lpv1.Answer_OPTION_B:
		if req.GetOptionB() == "" {
			return status.Error(codes.InvalidArgument, "option B must be provided if it is selected as the answer")
		}
	case lpv1.Answer_OPTION_C:
		if req.GetOptionC() == "" {
			return status.Error(codes.InvalidArgument, "option C must be provided if it is selected as the answer")
		}
	case lpv1.Answer_OPTION_D:
		if req.GetOptionD() == "" {
			return status.Error(codes.InvalidArgument, "option D must be provided if it is selected as the answer")
		}
	case lpv1.Answer_OPTION_E:
		if req.GetOptionE() == "" {
			return status.Error(codes.InvalidArgument, "option E must be provided if it is selected as the answer")
		}
	}

	return nil
}
