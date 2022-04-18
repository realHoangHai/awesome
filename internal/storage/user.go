package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/realHoangHai/awesome/internal/biz"
	"github.com/realHoangHai/awesome/internal/storage/ent"
	"github.com/realHoangHai/awesome/internal/storage/ent/user"
	"github.com/realHoangHai/awesome/pkg/log"
	"github.com/realHoangHai/awesome/pkg/utils"
	"time"
)

var _ biz.UserRepo = (*userRepo)(nil)

var userCacheKey = func(username string) string {
	return "user_cache_key" + username
}

type userRepo struct {
	store *Store
}

func NewUserRepo(store *Store) biz.UserRepo {
	return &userRepo{store: store}
}

func (r *userRepo) CreateUser(ctx context.Context, arg *biz.User) (*biz.User, error) {
	ph, err := utils.HashPassword(arg.Password)
	if err != nil {
		return nil, err
	}
	result, err := r.store.db.User.Create().SetUsername(arg.Username).SetPasswordHash(ph).Save(ctx)
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Id:       result.ID,
		Username: result.Username,
	}, nil
}

func (r *userRepo) GetUser(ctx context.Context, id int64) (*biz.User, error) {
	// try to fetch from cache
	cacheKey := userCacheKey(fmt.Sprintf("%d", id))
	target, err := r.getUserFromCache(ctx, cacheKey)
	if err != nil {
		// fetch from db while cache missed
		target, err = r.store.db.User.Get(ctx, id)
		if err != nil {
			return nil, biz.ErrUserNotFound
		}
		// set cache
		r.setUserCache(ctx, target, cacheKey)
	}
	return &biz.User{Id: target.ID, Username: target.Username}, nil
}

func (r *userRepo) VerifyPassword(ctx context.Context, u *biz.User) (bool, error) {
	result, err := r.store.db.User.Query().Where(user.UsernameEQ(u.Username)).Only(ctx)
	if err != nil {
		return false, err
	}
	return utils.CheckPassword(u.Password, result.PasswordHash), nil
}

func (r *userRepo) FindByUsername(ctx context.Context, username string) (*biz.User, error) {
	var target *ent.User
	// try to fetch from cache
	cacheKey := userCacheKey(username)
	target, err := r.getUserFromCache(ctx, cacheKey)
	if err != nil {
		// fetch from db while cache missed\
		target, err := r.store.db.User.
			Query().
			Where(user.UsernameEQ(username)).
			Only(ctx)
		if err != nil {
			return nil, biz.ErrUserNotFound
		}
		// set cache
		r.setUserCache(ctx, target, cacheKey)
	}
	return &biz.User{
		Id:       target.ID,
		Username: target.Username,
		Password: target.PasswordHash,
	}, nil
}

func (r *userRepo) getUserFromCache(ctx context.Context, key string) (*ent.User, error) {
	result, err := r.store.redisCli.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var cacheUser = &ent.User{}
	if err := json.Unmarshal([]byte(result), cacheUser); err != nil {
		return nil, err
	}
	return cacheUser, nil
}

func (r *userRepo) setUserCache(ctx context.Context, user *ent.User, key string) {
	marshal, err := json.Marshal(user)
	if err != nil {
		log.Errorf("fail to set user cache:json.Marshal(%v) error(%v)", user, err)
	}
	if err := r.store.redisCli.Set(ctx, key, string(marshal), time.Minute*30).Err(); err != nil {
		log.Errorf("fail to set user cache:redis.Set(%v) error(%v)", user, err)
	}
}
