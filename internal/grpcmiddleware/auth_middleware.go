package grpcmiddleware

import (
	"context"

	jwtEntity "github.com/irvandMR/go-grpc-ecommerce-BE/internal/entity/jwt"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/utils"
	gocache "github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
)



type authMiddleware struct {
	cacheService *gocache.Cache
}

func (am *authMiddleware) Middleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	if info.FullMethod == "/auth.AuthService/Login" || info.FullMethod == "/auth.AuthService/Register" {
		return handler(ctx, req)
	}

	//  Get token from metadata
	tknStr, err := jwtEntity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	_, ok := am.cacheService.Get(tknStr)
	if ok {
		return nil, utils.UnauthenticatedResponse()
	}

	// parse token to entity jwt
	claim , err := jwtEntity.GetClaimsFromToken(tknStr)
	if err != nil {
		return nil, err
	}	

	// entity jwt to context
	ctx = claim.SteToContext(ctx)

	res, err := handler(ctx, req)

	return res, err
}

func NewAuthMiddleware(cacheService *gocache.Cache) *authMiddleware {
	return &authMiddleware{
		cacheService: cacheService,
	}
}