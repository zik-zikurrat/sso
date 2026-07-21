package cache

import (
	"context"
	"log/slog"
	"sso/internal/config"
	"sso/internal/helper"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Redis struct {
	log    *slog.Logger
	Client *goredis.Client
}

func New(cfg *config.Config, l *slog.Logger) (*Redis, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:         helper.GetRedisAddr(cfg),
		DB:           cfg.Cache.DB,
		DialTimeout:  5 * time.Second, // ожидание установления tcp соединения
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		l.Error("failed to ping redis", slog.String("error", err.Error()))
		return nil, err
	}
	redis := &Redis{
		log:    l,
		Client: client,
	}
	return redis, nil
}

func (r *Redis) Close() {
	if r.Client != nil {
		err := r.Client.Close()
		if err != nil {
			r.log.Error("failed to close redis client", slog.String("error", err.Error()))
		}
	}
}
