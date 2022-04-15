package biz

import (
	"context"
)

type Card struct {
	Id      int64
	CardNo  string
	CCV     string
	Expires string
}

type CardRepo interface {
	CreateCard(ctx context.Context, c *Card) (*Card, error)
	GetCard(ctx context.Context, id int64) (*Card, error)
	ListCard(ctx context.Context, uid int64) ([]*Card, error)
}

type CardBiz struct {
	repo CardRepo
}

func NewCardBiz(repo CardRepo) *CardBiz {
	return &CardBiz{repo: repo}
}

func (biz *CardBiz) Create(ctx context.Context, c *Card) (*Card, error) {
	return biz.repo.CreateCard(ctx, c)
}

func (biz *CardBiz) Get(ctx context.Context, id int64) (*Card, error) {
	return biz.repo.GetCard(ctx, id)
}

func (biz *CardBiz) List(ctx context.Context, uid int64) ([]*Card, error) {
	return biz.repo.ListCard(ctx, uid)
}
