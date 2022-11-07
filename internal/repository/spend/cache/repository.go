package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/redis"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/utils"
)

type SpendRepo interface {
	GetByTimeSince(ctx context.Context, userId int64, timeSince time.Time) ([]model.Spend, error)
	Save(ctx context.Context, userId int64, timeSince time.Time, spends []model.Spend) error
	DeleteForUser(ctx context.Context, userId int64) error
}

type spendRedisRepo struct {
	redisClient redis.Client
}

func New(redisClient redis.Client) SpendRepo {
	return &spendRedisRepo{
		redisClient: redisClient,
	}
}

func (s *spendRedisRepo) GetByTimeSince(ctx context.Context, userId int64, timeSince time.Time) ([]model.Spend, error) {
	cacheKey := getCacheKey(userId, timeSince)
	res, err := s.redisClient.Get(ctx, cacheKey)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res == "" {
		return nil, nil
	}

	var spends []model.Spend
	err = msgpack.Unmarshal([]byte(res), &spends)
	if err != nil {
		return spends, errors.WithStack(err)
	}

	return spends, nil
}

func (s *spendRedisRepo) DeleteForUser(ctx context.Context, userId int64) error {
	now := time.Now()
	return s.redisClient.Del(ctx,
		getCacheKey(userId, utils.BeginOfDay(now)),
		getCacheKey(userId, utils.BeginOfMonth(now)),
		getCacheKey(userId, utils.BeginOfYear(now)),
	)
}

func (s *spendRedisRepo) Save(ctx context.Context, userId int64, timeSince time.Time, spends []model.Spend) error {
	cacheValue, err := msgpack.Marshal(spends)
	if err != nil {
		return err
	}

	return s.redisClient.Set(ctx, getCacheKey(userId, timeSince), cacheValue, time.Hour*24)
}

func getCacheKey(userId int64, timeSince time.Time) string {
	return fmt.Sprintf("spend_%d_%d", userId, timeSince.Unix())
}
