package repo

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/realHoangHai/awesome/config"
	"github.com/realHoangHai/awesome/internal/repo/ent"
	"github.com/realHoangHai/awesome/internal/repo/ent/migrate"
	"github.com/realHoangHai/awesome/pkg/log"
	"time"
	// init mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// ProviderSet is repository providers.
var ProviderSet = wire.NewSet(NewStore, NewEntClient, NewRedisCmd, NewUserRepo, NewCardRepo, NewAddressRepo)

// Store .
type Store struct {
	db       *ent.Client
	redisCli redis.Cmdable
}

func NewEntClient(cfg *config.Config, logger log.Logger) *ent.Client {
	client, err := ent.Open(
		cfg.DB.Driver,
		cfg.DB.Source,
	)
	if err != nil {
		logger.Fatalf("failed opening connection to db: %v", err)
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background(), migrate.WithForeignKeys(false)); err != nil {
		logger.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}

func NewRedisCmd(cfg *config.Config, logger log.Logger) redis.Cmdable {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
		DialTimeout:  time.Second * 2,
		PoolSize:     10,
	})
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	err := client.Ping(timeout).Err()
	if err != nil {
		logger.Fatalf("redis connect error: %v", err)
	}
	return client
}

// NewStore .
func NewStore(entClient *ent.Client, redisCmd redis.Cmdable, logger log.Logger) (*Store, func(), error) {
	store := &Store{
		db:       entClient,
		redisCli: redisCmd,
	}
	return store, func() {
		if err := store.db.Close(); err != nil {
			logger.Error(err)
		}
	}, nil
}
