package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/handler"
	"github.com/irvandMR/go-grpc-ecommerce-BE/pb/service"
	"github.com/irvandMR/go-grpc-ecommerce-BE/pkg/database"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func errorMidleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res any, err error) {
	res, err = handler(ctx, req)

	
}

func main() {

	ctx := context.Background()

	godotenv.Load()
	listen, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Panicf("Error starting TCP server: %v", err)
	}
	database.ConnectDB(ctx, os.Getenv("DB_URI"))
	serviceHandler := handler.NewServiceHandler()
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor()
	)
	service.RegisterHelloWorldServiceServer(grpcServer, serviceHandler)

	if os.Getenv("ENVIROMENT") == "dev" {
		reflection.Register(grpcServer)
		log.Println("reflection service registered")
	}

	

	log.Println("Starting gRPC server on :", listen.Addr().String(), "in", os.Getenv("ENVIROMENT"), "mode")
	if err := grpcServer.Serve(listen); err != nil {
		log.Panicf("Error starting gRPC server: %v", err)
	}
}