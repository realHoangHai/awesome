package repo

import (
	"context"
	"github.com/realHoangHai/awesome/internal/biz"
	"github.com/realHoangHai/awesome/pkg/log"
)

var _ biz.AddressRepo = (*addressRepo)(nil)

type addressRepo struct {
	store *Store
	log   *log.Logh
}

func NewAddressRepo(store *Store, log log.Logger) biz.AddressRepo {
	return &addressRepo{
		store: store,
		log:   nil,
	}
}

func (r *addressRepo) CreateAddress(ctx context.Context, a *biz.Address) (*biz.Address, error) {
	result, err := r.store.db.Address.
		Create().
		SetName(a.Name).
		SetAddress(a.Address).
		SetMobile(a.Mobile).
		SetPostCode(a.PostCode).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &biz.Address{
		Id:       result.ID,
		Name:     result.Name,
		Mobile:   result.Mobile,
		Address:  result.Address,
		PostCode: result.PostCode,
	}, nil
}

func (r *addressRepo) GetAddress(ctx context.Context, id int64) (*biz.Address, error) {
	result, err := r.store.db.Address.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &biz.Address{
		Id:       result.ID,
		Name:     result.Name,
		Mobile:   result.Mobile,
		PostCode: result.PostCode,
		Address:  result.Address,
	}, nil
}

func (r *addressRepo) ListAddress(ctx context.Context, uid int64) ([]*biz.Address, error) {
	list, err := r.store.db.Address.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*biz.Address, 0)
	for _, addr := range list {
		result = append(result, &biz.Address{
			Id:       addr.ID,
			Name:     addr.Name,
			Mobile:   addr.Mobile,
			PostCode: addr.PostCode,
			Address:  addr.Address,
		})
	}
	return result, nil
}
