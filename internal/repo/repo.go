package repo

import (
	"context"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/realHoangHai/awesome/internal/repo/ent"
	"github.com/realHoangHai/awesome/internal/repo/ent/migrate"
	"time"

	// init mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// ProviderSet is repository providers.
var ProviderSet = wire.NewSet()

// Store .
type Store struct {
	db       *ent.Client
	redisCli redis.Cmdable
}

func NewEntClient(cfg *config.Config) *ent.Client {
	client, err := ent.Open(
		cfg.Database.Driver,
		cfg.Database.Source,
	)
	if err != nil {
		log.Fatalf("failed opening connection to db: %v", err)
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background(), migrate.WithForeignKeys(false)); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}

func NewRedisCmd(cfg *config.Config) redis.Cmdable {
	//log := log.NewHelper(log.With(logger, "module", "user-service/data/ent"))
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		ReadTimeout:  cfg.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: cfg.Redis.WriteTimeout.AsDuration(),
		DialTimeout:  time.Second * 2,
		PoolSize:     10,
	})
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	err := client.Ping(timeout).Err()
	if err != nil {
		log.Fatalf("redis connect error: %v", err)
	}
	return client
}

// NewStore .
func NewStore(entClient *ent.Client, redisCmd redis.Cmdable) (*Store, func(), error) {
	//log := log.NewHelper(log.With(logger, "module", "user-service/data"))

	d := &Store{
		db:       entClient,
		redisCli: redisCmd,
	}
	return d, func() {
		if err := d.db.Close(); err != nil {
			log.Error(err)
		}
	}, nil
}
