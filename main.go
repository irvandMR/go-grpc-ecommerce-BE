package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/handler"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/repository"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/service"
	"github.com/irvandMR/go-grpc-ecommerce-BE/pb/auth"
	"github.com/irvandMR/go-grpc-ecommerce-BE/pkg/database"
	"github.com/irvandMR/go-grpc-ecommerce-BE/pkg/grpcmiddleware"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)


func main() {

	ctx := context.Background()

	godotenv.Load()
	listen, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Panicf("Error starting TCP server: %v", err)
	}
	db := database.ConnectDB(ctx, os.Getenv("DB_URI"))

	// Auth
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpcmiddleware.ErrorMiddleware),
	)

	// Register your gRPC services here
	auth.RegisterAuthServiceServer(grpcServer, authHandler)

	if os.Getenv("ENVIROMENT") == "dev" {
		reflection.Register(grpcServer)
		log.Println("reflection service registered")
	}

	

	log.Println("Starting gRPC server on :", listen.Addr().String(), "in", os.Getenv("ENVIROMENT"), "mode")
	if err := grpcServer.Serve(listen); err != nil {
		log.Panicf("Error starting gRPC server: %v", err)
	}
}