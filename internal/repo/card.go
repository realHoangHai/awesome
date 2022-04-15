package repo

import (
	"context"
	"github.com/realHoangHai/awesome/internal/biz"
	"github.com/realHoangHai/awesome/internal/repo/ent/user"
)

var _ biz.CardRepo = (*cardRepo)(nil)

type cardRepo struct {
	store *Store
}

func NewCardRepo(store *Store) biz.CardRepo {
	return &cardRepo{store: store}
}

func (r *cardRepo) CreateCard(ctx context.Context, c *biz.Card) (*biz.Card, error) {
	result, err := r.store.db.Card.
		Create().
		SetCardNo(c.CardNo).
		SetCcv(c.CCV).
		SetExpires(c.Expires).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &biz.Card{
		Id:      result.ID,
		CardNo:  result.CardNo,
		CCV:     result.Ccv,
		Expires: result.Expires,
	}, nil
}

func (r *cardRepo) GetCard(ctx context.Context, id int64) (*biz.Card, error) {
	result, err := r.store.db.Card.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &biz.Card{
		Id:      result.ID,
		CardNo:  result.CardNo,
		CCV:     result.Ccv,
		Expires: result.Expires,
	}, nil
}

func (r *cardRepo) ListCard(ctx context.Context, uid int64) ([]*biz.Card, error) {
	list, err := r.store.db.User.
		Query().
		Where(user.ID(uid)).
		QueryCards().
		All(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*biz.Card, 0)
	for _, c := range list {
		result = append(result, &biz.Card{
			Id:      c.ID,
			CardNo:  c.CardNo,
			CCV:     c.Ccv,
			Expires: c.Expires,
		})
	}
	return result, nil
}
