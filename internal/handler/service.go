package handler

import (
	"context"
	"fmt"

	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/utils"
	"github.com/irvandMR/go-grpc-ecommerce-BE/pb/service"
)

type serviceHandler struct {
	service.UnimplementedHelloWorldServiceServer
}

func (sh * serviceHandler) HelloWorld(ctx context.Context,req *service.HelloWorldRequest) (*service.HelloWorldResponse, error) {
	errRes, err := utils.CheckValidtion(req)
	if err != nil {
		return nil, err
	}
	if errRes != nil {
		return &service.HelloWorldResponse{
			BaseResponse: utils.ErrorResponse(errRes),
		}, nil
	}	
	return &service.HelloWorldResponse{
		Message: fmt.Sprintf("Hello %s, Welcome to gRPC World!", req.Name),
		BaseResponse: utils.SuccessResponse("Success"),
	}, nil
}

func NewServiceHandler() *serviceHandler {
	return &serviceHandler{}
} 