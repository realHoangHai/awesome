package service

import (
	"context"
	v1 "github.com/realHoangHai/awesome/api/user/v1"
	"github.com/realHoangHai/awesome/internal/biz"
	"github.com/realHoangHai/awesome/pkg/log"
)

type UserService struct {
	v1.UnimplementedUserServer

	ub  *biz.UserBiz
	ab  *biz.AddressBiz
	cb  *biz.CardBiz
	log *log.Logh
}

func NewUserService(ub *biz.UserBiz, ab *biz.AddressBiz, cb *biz.CardBiz, logger log.Logger) *UserService {
	return &UserService{
		ub:  ub,
		ab:  ab,
		cb:  cb,
		log: nil,
	}
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
