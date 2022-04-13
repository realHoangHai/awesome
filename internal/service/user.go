package service

import (
	"context"
	v1 "github.com/realHoangHai/awesome/api/user/v1"
	"github.com/realHoangHai/awesome/internal/biz"
	"github.com/realHoangHai/awesome/pkg/log"
)

type UserService struct {
	v1.UnimplementedUserServer
	biz *biz.UserBiz
	log log.Logger
}

func NewUserService(ub *biz.UserBiz, logger log.Logger) *UserService {
	return &UserService{
		biz: ub,
		log: logger,
	}
}

func (s *UserService) CreateUser(ctx context.Context, in *v1.SaveUserReq) (*v1.SaveUserReply, error) {
	return s.biz.Save(ctx, in)
}

func (s *UserService) GetUser(ctx context.Context, req *v1.GetUserReq) (*v1.GetUserReply, error) {
	result, err := s.biz.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetUserReply{
		Id:       result.Id,
		Username: result.Username,
	}, nil
}

func (s *UserService) VerifyPassword(ctx context.Context, req *v1.VerifyPasswordReq) (*v1.VerifyPasswordReply, error) {
	result, err := s.biz.VerifyPassword(ctx, &biz.User{Username: req.Username, Password: req.Password})
	if err != nil {
		return nil, err
	}
	return &v1.VerifyPasswordReply{
		Ok: result,
	}, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, in *v1.GetUserByUsernameReq) (*v1.GetUserByUsernameReply, error) {
	return s.biz.GetUserByUserName(ctx, in)
}
