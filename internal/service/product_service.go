package service

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/entity"
	jwtEntity "github.com/irvandMR/go-grpc-ecommerce-BE/internal/entity/jwt"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/repository"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/utils"
	"github.com/irvandMR/go-grpc-ecommerce-BE/pb/product"
)

type IProductService interface {
	CreateProduct(ctx context.Context, req *product.CreateProductReqauest) (*product.CreateProductResponse, error)
}

type productService struct {
	productRepo     repository.IProductRepository
}

func (ps *productService) CreateProduct(ctx context.Context, req *product.CreateProductReqauest) (*product.CreateProductResponse, error) {

	// check user is admin or not
	ctxUser, err := jwtEntity.GetClaimsFromContext(ctx)
	if err != nil{
		return nil, err
	}
	if(ctxUser.Role != entity.ADMIN_ROLE){
		return  nil,utils.UnauthenticatedResponse()
	}
	// check already have image or not
	imagePath := filepath.Join("storage", "product_images", req.ImageFileName)
	_, err = os.Stat(imagePath)
	if err != nil{
		if os.IsNotExist(err){

			return &product.CreateProductResponse{
				Base: utils.BadRequestResponse("image file not found"),
			}, nil
		}
		return nil, err
	}
	
	// insert to db
	// Process Create New Product logic here
	productEntity := entity.Product{
		Id: uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ImageFileName: req.ImageFileName,
		CreatedAt: time.Now(),
		CreatedBy: ctxUser.Fullname,
	}
	err = ps.productRepo.CreateNewProduct(ctx, &productEntity);
	if err != nil {
		return nil, err
	}
	return &product.CreateProductResponse{
		Base: utils.SuccessResponse("success create new product"),
		Id: productEntity.Id,
	}, nil;
}

func NewProductService(productRepo repository.IProductRepository) IProductService {
	return &productService{
		productRepo:     productRepo,
	}
}