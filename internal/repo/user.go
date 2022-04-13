package repo

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/realHoangHai/awesome/internal/biz"
)

var _ biz.UserRepo = (*UserRepo)(nil)

var userCacheKey = func(username string) string {
	return "user_cache_key" + username
}

type userRepo struct {
	store *Store
	log   *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/server-service")),
	}
}
