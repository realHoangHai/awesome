package biz

import (
	"context"
	"errors"
	v1 "github.com/realHoangHai/awesome/api/user/v1"
	"github.com/realHoangHai/awesome/pkg/log"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	Id       int64
	Username string
	Password string
}

type UserRepo interface {
	CreateUser(ctx context.Context, u *User) (*User, error)
	GetUser(ctx context.Context, id int64) (*User, error)
	VerifyPassword(ctx context.Context, u *User) (bool, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
}

type UserBiz struct {
	repo UserRepo
}

func NewUserBiz(repo UserRepo) *UserBiz {
	return &UserBiz{repo: repo}
}

func (biz *UserBiz) CreateUser(ctx context.Context, u *User) (*User, error) {
	result, err := biz.repo.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (biz *UserBiz) GetUser(ctx context.Context, id int64) (*User, error) {
	return biz.repo.GetUser(ctx, id)
}

func (biz *UserBiz) VerifyPassword(ctx context.Context, u *User) (bool, error) {
	return biz.repo.VerifyPassword(ctx, u)
}

func (biz *UserBiz) GetUserByUserName(ctx context.Context, req *v1.GetUserByUsernameReq) (*v1.GetUserByUsernameReply, error) {
	user, err := biz.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		log.Errorf("get user by username error: %v", err)
		return nil, err
	}
	return &v1.GetUserByUsernameReply{
		Id:       user.Id,
		Username: user.Username,
	}, nil
}
