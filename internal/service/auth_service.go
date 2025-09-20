package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/entity"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/repository"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/utils"
	"github.com/irvandMR/go-grpc-ecommerce-BE/pb/auth"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IAuthService interface {
	Register(ctx context.Context, req *auth.RegisterRequest) (res *auth.RegisterResponse, err error)
	Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error)
}
type authService struct{
	authRepo repository.IAuthRepository
}

// Login implements the IAuthService interface.
func (au *authService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	// check email to db
	user, err := au.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil{
		return nil, err
	}

	if user == nil{
		return &auth.LoginResponse{
			Base: utils.BadRequestResponse("email not registered"),
		}, nil
	}

	// if exist, check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil{
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword){
			return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
		}
		return nil, err
	}
	// if match, generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.JWTClaim{
		Email: user.Email,
		Fullname: user.Fullname,
		Role: user.RoleCode,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: user.Id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	})
	secretKey := os.Getenv("JWT_SECRET")  // should be from env
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil{
		return nil, err
	}

	// send token to user

	return &auth.LoginResponse{
		Base: utils.SuccessResponse("Success Login"),
		Token: signedToken,
	}, nil
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