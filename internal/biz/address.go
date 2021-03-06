package biz

import (
	"context"
)

type Address struct {
	Id       int64
	Name     string
	Mobile   string
	Address  string
	PostCode string
}

type AddressRepo interface {
	CreateAddress(ctx context.Context, a *Address) (*Address, error)
	GetAddress(ctx context.Context, id int64) (*Address, error)
	ListAddress(ctx context.Context, uid int64) ([]*Address, error)
}

type AddressBiz struct {
	repo AddressRepo
}

func NewAddressBiz(repo AddressRepo) *AddressBiz {
	return &AddressBiz{repo: repo}
}

func (biz *AddressBiz) Create(ctx context.Context, uid int64, a *Address) (*Address, error) {
	return biz.repo.CreateAddress(ctx, a)
}

func (biz *AddressBiz) Get(ctx context.Context, id int64) (*Address, error) {
	return biz.repo.GetAddress(ctx, id)
}

func (biz *AddressBiz) List(ctx context.Context, uid int64) ([]*Address, error) {
	return biz.repo.ListAddress(ctx, uid)
}
