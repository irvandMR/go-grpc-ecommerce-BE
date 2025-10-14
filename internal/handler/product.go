package handler

import (
	"context"

	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/service"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/utils"
	"github.com/irvandMR/go-grpc-ecommerce-BE/pb/product"
)

type productHandler struct {
	product.UnimplementedProductServiceServer
	productService service.IProductService
}

func (ph *productHandler) CreateProduct(ctx context.Context, req *product.CreateProductReqauest) (*product.CreateProductResponse, error) {
		errRes, err := utils.CheckValidtion(req)
		if err != nil {
			return nil, err
		}
		if errRes != nil {
			return &product.CreateProductResponse{
				Base: utils.ErrorResponse(errRes),
			}, nil
		}
		// Process Register logic here
		res, errAuth := ph.productService.CreateProduct(ctx, req)
		if errAuth != nil {
			return nil, errAuth
		}
		return res, nil
}

func NewProductHandler(productService service.IProductService) *productHandler {
	return &productHandler{
		productService: productService,
	}
} 