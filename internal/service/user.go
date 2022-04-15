package service

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/realHoangHai/awesome/api/user/v1"
	"github.com/realHoangHai/awesome/internal/biz"
	"google.golang.org/grpc"
)

type UserService struct {
	v1.UnimplementedUserServer

	ub *biz.UserBiz
	ab *biz.AddressBiz
	cb *biz.CardBiz
}

func NewUserService(ub *biz.UserBiz, ab *biz.AddressBiz, cb *biz.CardBiz) *UserService {
	return &UserService{ub: ub, ab: ab, cb: cb}
}

// Register implements server.Service interface.
func (s *UserService) Register(srv *grpc.Server) {
	v1.RegisterUserServer(srv, s)
}

// RegisterWithEndpoint implements server.EndpointService interface.
func (s *UserService) RegisterWithEndpoint(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) {
	_ = v1.RegisterUserHandlerFromEndpoint(ctx, mux, addr, opts)
}

func (s *UserService) CreateUser(ctx context.Context, req *v1.CreateUserReq) (*v1.CreateUserReply, error) {
	result, err := s.ub.CreateUser(ctx, &biz.User{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateUserReply{
		Id:       result.Id,
		Username: result.Username,
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *v1.GetUserReq) (*v1.GetUserReply, error) {
	result, err := s.ub.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserReply{
		Id:       result.Id,
		Username: result.Username,
	}, nil
}

func (s *UserService) VerifyPassword(ctx context.Context, req *v1.VerifyPasswordReq) (*v1.VerifyPasswordReply, error) {
	result, err := s.ub.VerifyPassword(ctx, &biz.User{Username: req.Username, Password: req.Password})
	if err != nil {
		return nil, err
	}
	return &v1.VerifyPasswordReply{
		Ok: result,
	}, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, req *v1.GetUserByUsernameReq) (*v1.GetUserByUsernameReply, error) {
	return s.ub.GetUserByUserName(ctx, req)
}
