package jwt

import (
	"context"
	"strings"

	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/utils"
	"google.golang.org/grpc/metadata"
)

func ParseTokenFromContext(ctx context.Context) (string, error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", utils.UnauthenticatedResponse()
	}
	bearerToken, ok := metadata["authorization"]
	if !ok || len(bearerToken) == 0 {
		return "", utils.UnauthenticatedResponse()
	}

	tokensSplit := strings.Split(bearerToken[0], " ")
	if len(tokensSplit) != 2 || strings.ToLower(tokensSplit[0]) != "bearer" {
		return "", utils.UnauthenticatedResponse()
	}
	jwtToken := tokensSplit[1]
	return jwtToken, nil

}