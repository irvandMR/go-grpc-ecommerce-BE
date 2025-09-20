package handler

import (
	"context"

	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/service"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/utils"
	"github.com/irvandMR/go-grpc-ecommerce-BE/pb/auth"
)

type authHandler struct {
	auth.UnimplementedAuthServiceServer
	authService service.IAuthService
}

func (au *authHandler) Register(ctx context.Context, req *auth.RegisterRequest) (res *auth.RegisterResponse, err error) {
		errRes, err := utils.CheckValidtion(req)
		if err != nil {
			return nil, err
		}
		if errRes != nil {
			return &auth.RegisterResponse{
				Base: utils.ErrorResponse(errRes),
			}, nil
		}
		// Process Register logic here
		res, errAuth := au.authService.Register(ctx, req)
		if errAuth != nil {
			return nil, errAuth
		}
		return res, nil
}



func NewAuthHandler(authService service.IAuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
} 