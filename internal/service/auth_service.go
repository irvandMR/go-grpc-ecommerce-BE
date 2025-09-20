package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/entity"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/repository"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/utils"
	"github.com/irvandMR/go-grpc-ecommerce-BE/pb/auth"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Register(ctx context.Context, req *auth.RegisterRequest) (res *auth.RegisterResponse, err error)
}
type authService struct{
	authRepo repository.IAuthRepository
}

func (au *authService) Register(ctx context.Context, req *auth.RegisterRequest) (res *auth.RegisterResponse, err error) {

	if req.Password != req.PasswordConfirmation{
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("password and confirm password not match"),
		}, nil
	}
	// Check email to db
	user, err := au.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil{
		return nil, err
	}
	// if email exist already, we void error
	if user != nil{
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("email already registered"),
		}, nil
	}
	// if not exist, we hash password and save to db
	hashPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil{
		return nil, err
	}

	// insert to db
	newUser := &entity.User{
		Id:  uuid.NewString(),
		Email: req.Email,
		Password: string(hashPass),
		Fullname: req.FullName,
		RoleCode: entity.CUSTOMER_ROLE,
		CreatedAt: time.Now(),
		CreatedBy: &req.FullName,
	}
	errAuthIns := au.authRepo.InsertUser(ctx, newUser)
	if errAuthIns != nil{
		return nil, errAuthIns
	}


	return &auth.RegisterResponse{
		Base: utils.SuccessResponse("Success Register User"),
	}, nil
}



func NewAuthService(authRepo repository.IAuthRepository) IAuthService {
	return &authService{
		authRepo: authRepo,
	}
}