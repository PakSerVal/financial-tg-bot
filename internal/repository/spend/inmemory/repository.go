package inmemory

import (
	"context"
	"database/sql"
	"time"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
)

type inmemory struct {
	lastIndex int64
	records   map[int64][]model.Spend
}

func New() spend.Repository {
	return &inmemory{
		records: map[int64][]model.Spend{},
	}
}

func (i *inmemory) SaveTx(tx *sql.Tx, ctx context.Context, sum int64, category string, userId int64) error {
	rec := model.Spend{
		Id:        i.lastIndex + 1,
		Price:     sum,
		Category:  category,
		CreatedAt: time.Now(),
		UserId:    userId,
	}
	i.records[userId] = append(i.records[userId], rec)

	return nil
}

func (i *inmemory) GetByTimeSinceTx(tx *sql.Tx, ctx context.Context, userId int64, timeSince time.Time) ([]model.Spend, error) {
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
