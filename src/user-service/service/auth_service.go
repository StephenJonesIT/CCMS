package service

import (
	"context"
	"errors"	
	"time"
	"github.com/StephenJonesIT/CCMS/src/user-service/common"
	"github.com/StephenJonesIT/CCMS/src/user-service/internal/repository"
	pb "github.com/StephenJonesIT/CCMS/src/user-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AuthService struct {
    pb.UnimplementedAuthServiceServer
    userRepo    repository.UserRepository
    jwtSecret   string
    tokenExpiry time.Duration
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, tokenExpiry time.Duration) *AuthService {
    return &AuthService{
        userRepo:    userRepo,
        jwtSecret:   jwtSecret,
        tokenExpiry: tokenExpiry,
    }
}

func (s *AuthService) VerifyToken(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
    claims, err := common.VerifyToken(req.Token)
    if err != nil {
        return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
    }

    // Verify user exists
    _, err = s.userRepo.GetUserByID(claims.UserID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, status.Error(codes.NotFound, "user not found")
        }
        return nil, status.Error(codes.Internal, "database error")
    }

    return &pb.TokenResponse{
        UserId:   claims.UserID.String(),
        RoleName: claims.RoleName,
    }, nil
}

func (s *AuthService) CheckPermission(ctx context.Context, req *pb.PermissionRequest) (*pb.PermissionResponse, error) {
    claims, err := common.VerifyToken(req.Token)
    if err != nil {
        return nil, status.Error(codes.Unauthenticated, "invalid token")
    }

    // Check permission
    requiredPermission := req.Resource
    hasPerm := s.userRepo.HasPermission(claims.UserID, requiredPermission)
    
    return &pb.PermissionResponse{
        Allowed:  hasPerm,
        RoleName: claims.RoleName,
    }, nil
}
