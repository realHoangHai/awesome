package biz

import (
	"context"
	v1 "github.com/realHoangHai/awesome/api/user/v1"
	"github.com/realHoangHai/awesome/pkg/log"
	"math/rand"
)

type User struct {
	Id       int64
	Username string
	Password string
}

type UserRepo interface {
	CreateUser(context.Context, *User) (*User, error)
	GetUser(context.Context, int64) (*User, error)
	VerifyPassword(context.Context, *User) (bool, error)
	FindByCondition(context.Context, map[string]interface{}) (*User, error)
}

type UserBiz struct {
	repo UserRepo
	log  log.Logger
}

func NewUserBiz(repo UserRepo, logger log.Logger) *UserBiz {
	return &UserBiz{repo: repo, log: logger}
}

func (biz *UserBiz) Save(ctx context.Context, in *v1.SaveUserReq) (*v1.SaveUserReply, error) {
	user := &User{
		Id:       rand.Int63(),
		Username: in.Username,
		Password: in.Password,
	}
	_, err := biz.repo.CreateUser(ctx, user)
	if err != nil {
		biz.log.Errorf("save user error: %v", err)
		return nil, err
	}
	return &v1.SaveUserReply{
		Id: user.Id,
	}, nil
}

func (biz *UserBiz) GetUser(ctx context.Context, id int64) (*User, error) {
	return biz.repo.GetUser(ctx, id)
}

func (biz *UserBiz) VerifyPassword(ctx context.Context, u *User) (bool, error) {
	return biz.repo.VerifyPassword(ctx, u)
}

func (biz *UserBiz) GetUserByUserName(ctx context.Context, in *v1.GetUserByUsernameReq) (*v1.GetUserByUsernameReply, error) {
	user, err := biz.repo.FindByCondition(ctx, map[string]interface{}{"username": in.Username})
	if err != nil {
		biz.log.Errorf("get user by username error: %v", err)
		return nil, err
	}
	return &v1.GetUserByUsernameReply{
		Id:       user.Id,
		Username: user.Username,
	}, nil
}
