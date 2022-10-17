package inmemory

import (
	"context"
	"time"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
)

type inmemory struct {
	lastIndex int64
	records   map[int64][]model.Spend
}

func New() *inmemory {
	return &inmemory{
		records: map[int64][]model.Spend{},
	}
}

func (i *inmemory) Save(ctx context.Context, price int64, category string, userId int64) error {
	rec := model.Spend{
		Id:        i.lastIndex + 1,
		Price:     price,
		Category:  category,
		CreatedAt: time.Now(),
		UserId:    userId,
	}
	i.records[userId] = append(i.records[userId], rec)

	return nil
}

func (i *inmemory) GetByTimeSince(ctx context.Context, userId int64, timeSince time.Time) ([]model.Spend, error) {
	if _, ok := i.records[userId]; !ok {
		return nil, nil
	}

	var result []model.Spend
	for _, rec := range i.records[userId] {
		if timeSince.Before(rec.CreatedAt) {
			result = append(result, rec)
		}
	}

	return result, nil
}
