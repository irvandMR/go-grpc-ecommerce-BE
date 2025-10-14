package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/grpcmiddleware"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/handler"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/repository"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/service"
	"github.com/irvandMR/go-grpc-ecommerce-BE/pb/auth"
	"github.com/irvandMR/go-grpc-ecommerce-BE/pb/product"
	"github.com/irvandMR/go-grpc-ecommerce-BE/pkg/database"
	"github.com/joho/godotenv"
	gocache "github.com/patrickmn/go-cache"
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

	// Cache
	cacheService := gocache.New(time.Hour*24, time.Hour)

	authMiddleware :=  grpcmiddleware.NewAuthMiddleware(cacheService)

	// Auth
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, cacheService)
	authHandler := handler.NewAuthHandler(authService)

	// Product
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpcmiddleware.ErrorMiddleware, authMiddleware.Middleware),
	)

	// Register your gRPC services here
	auth.RegisterAuthServiceServer(grpcServer, authHandler)
	product.RegisterProductServiceServer(grpcServer, productHandler)

	if os.Getenv("ENVIROMENT") == "dev" {
		reflection.Register(grpcServer)
		log.Println("reflection service registered")
	}

	

	log.Println("Starting gRPC server on :", listen.Addr().String(), "in", os.Getenv("ENVIROMENT"), "mode")
	if err := grpcServer.Serve(listen); err != nil {
		log.Panicf("Error starting gRPC server: %v", err)
	}
}