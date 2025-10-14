package service

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/entity"
	jwtEntity "github.com/irvandMR/go-grpc-ecommerce-BE/internal/entity/jwt"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/repository"
	"github.com/irvandMR/go-grpc-ecommerce-BE/internal/utils"
	"github.com/irvandMR/go-grpc-ecommerce-BE/pb/auth"
	gocache "github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IAuthService interface {
	Register(ctx context.Context, req *auth.RegisterRequest) (res *auth.RegisterResponse, err error)
	Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error)
	Logout(ctx context.Context,req *auth.LogoutRequest) (*auth.LogoutResponse, error)
	ChangePassword(ctx context.Context,req *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error)
	GetProfile(ctx context.Context,req *auth.GetProfileRequest) (*auth.GetProfileResponse, error)
}
type authService struct{
	authRepo repository.IAuthRepository
	cacheService *gocache.Cache
}

// ChagePassword implements the IAuthService interface.
func (au *authService) ChangePassword(ctx context.Context, req *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	// check new password and confirm new password
	if req.NewPassword != req.NewPasswordConfirmation{
		return &auth.ChangePasswordResponse{
			Base: utils.BadRequestResponse("password and confirm password not match"),
		}, nil
	}
	// check current password with db
	jwtTkn, err := jwtEntity.ParseTokenFromContext(ctx)
	if err != nil{
		return nil, err
	}
	claims, err := jwtEntity.GetClaimsFromToken(jwtTkn)
	if err != nil{
		return nil, err
	}
	user, err := au.authRepo.GetUserByEmail(ctx, claims.Email)
	if err != nil{
		return nil, err
	}
	if user == nil{
		return &auth.ChangePasswordResponse{
			Base: utils.BadRequestResponse("email not registered"),
		}, nil
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
	if err != nil{
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword){
			return &auth.ChangePasswordResponse{
				Base: utils.BadRequestResponse("current password not match"),
			}, nil
		}
		return nil, err
	}

	// if match, update to new password
	hashPassword , err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
	if err != nil{
		return nil, err
	}
	err = au.authRepo.UpdatedUserPassword(ctx, user.Id, string(hashPassword), user.Fullname)
	if err != nil{
		return nil, err
	}

	return &auth.ChangePasswordResponse{
		Base: utils.SuccessResponse("Password changed successfully (not implemented)"),
	}, nil
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtEntity.JWTClaim{
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



func (au *authService) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	// get token from meta data
	jwtToken , err := jwtEntity.ParseTokenFromContext(ctx)
	if err != nil{
		return  nil, err
	}

	// returm token  to entity jwt
	claims, err := jwtEntity.GetClaimsFromContext(ctx)
	if err != nil{
		return nil,err
	}

	// instert token to cache
	au.cacheService.Set(jwtToken, "", time.Duration(claims.ExpiresAt.Time.Unix() - time.Now().Unix())*time.Second)

	return &auth.LogoutResponse{
		Base: utils.SuccessResponse("Success Logout"),
	}, nil
}

func (au *authService) GetProfile(ctx context.Context, req *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {
	// Get data token
	claims, err := jwtEntity.GetClaimsFromContext(ctx)
	if err != nil{
		return nil,err
	}
	
	
	// get data from db
	profile, err := au.authRepo.GetUserByEmail(ctx, claims.Email)
	if err != nil{
		return nil, err
	}
	if profile == nil{
		return &auth.GetProfileResponse{
			Base: utils.BadRequestResponse("email not registered"),
		
		}, nil
	}

	log.Println("claim", claims)

	return  &auth.GetProfileResponse{
		Base: utils.SuccessResponse("Success Get Profile"),
			Id:claims.Issuer,
			Email: claims.Email,
			FullName: claims.Fullname,
			RoleCode: claims.Role,
			MemberSince: timestamppb.New(profile.CreatedAt),
	},nil
}


func NewAuthService(authRepo repository.IAuthRepository, cacheService *gocache.Cache) IAuthService {
	return &authService{
		authRepo: authRepo,
		cacheService: cacheService,
	}
}