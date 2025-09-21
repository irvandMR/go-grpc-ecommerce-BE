package utils

import (
	"github.com/irvandMR/go-grpc-ecommerce-BE/pb/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SuccessResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatuseCode: 200,
		Message:     message,
	}
}

func ErrorResponse(validateErrors []*common.ValidationError) *common.BaseResponse {
	return &common.BaseResponse{
		StatuseCode:     400,
				Message: 	   "validation error",
				IsError:        true,
				ValidationErrors: validateErrors,
	}
}

func BadRequestResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatuseCode: 400,
		Message:     message,
		IsError:     true,
	}
}

func UnauthenticatedResponse() error {
	return status.Errorf(codes.Unauthenticated, "Unauthenticated")
}