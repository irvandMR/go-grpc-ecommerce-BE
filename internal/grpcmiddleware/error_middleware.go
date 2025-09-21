package grpcmiddleware

import (
	"context"
	"log"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "internal server error")
		}

	}()
	res, err := handler(ctx, req)
	if err != nil {
		log.Println(err)
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.Unauthenticated {
				return nil, err
			}
		}
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	return res, err
}