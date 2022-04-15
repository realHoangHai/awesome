package service

import (
	v1 "github.com/realHoangHai/awesome/api/user/v1"
	"github.com/realHoangHai/awesome/internal/biz"
)

type AddressService struct {
	v1.UnimplementedUserServer
	ab *biz.AddressBiz
}

func NewAddressService(biz *biz.AddressBiz) *UserService {
	return &UserService{ab: biz}
}
