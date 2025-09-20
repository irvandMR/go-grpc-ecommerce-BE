package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/irvandMR/go-grpc-ecommerce-BE/pb/service"
)

type serviceHandler struct {
	service.UnimplementedHelloWorldServiceServer
}

func (sh * serviceHandler) HelloWorld(ctx context.Context,req *service.HelloWorldRequest) (*service.HelloWorldResponse, error) {
	return &service.HelloWorldResponse{
		Message: fmt.Sprintf("Hello %s, Welcome to gRPC World!", req.Name),
	}, nil
}

func NewServiceHandler() *serviceHandler {
	return &serviceHandler{}
} 