package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/handler"
)

func main() {

	app := fiber.New()

	app.Use(cors.New())
	app.Get("/storage/product_images/:filename", handler.GetProductImageHandler)
	app.Post("/product/upload", handler.UploadProductImageHandler)

	app.Listen(":3000")
}