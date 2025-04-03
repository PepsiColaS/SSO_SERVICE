// Обработка хенделеров GRPC
package auth

import (
	"context"
	"sso/internal/database/postgres"

	ssov1 "github.com/PepsiColaS/SSO_PROTO/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, login string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, login string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	storage *postgres.PostgresClient
}

// Регистрируем обработчик запросов для gRPC сервера
func Register(gRPC *grpc.Server, storage *postgres.PostgresClient) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{storage: storage})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponce, error) {
	if req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if req.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "appID is required")
	}

	token, err := s.storage.Login(req.GetLogin(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponce{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponce, error) {
	if req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	userID, err := s.storage.RegisterNewUser(req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponce{
		UserId: userID,
	}, nil
}

func (s *serverAPI) isAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponce, error) {
	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "userID is required")
	}
	isAdmin, err := s.storage.IsAdmin(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.IsAdminResponce{
		IsAdmin: isAdmin,
	}, nil
}
