package service

import (
	v1 "github.com/realHoangHai/awesome/api/user/v1"
	"github.com/realHoangHai/awesome/internal/biz"
	"github.com/realHoangHai/awesome/pkg/log"
)

type AddressService struct {
	v1.UnimplementedUserServer
	biz *biz.UserBiz
	log *log.Logh
}

func NewAddressService(biz *biz.UserBiz, logger log.Logger) *UserService {
	return &UserService{
		biz: biz,
		log: nil,
	}
}
