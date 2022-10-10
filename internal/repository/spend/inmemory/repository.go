package inmemory

import (
	"time"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
)

type inmemory struct {
	lastIndex int64
	records   []model.Spend
}

func New() spend.Repository {
	return &inmemory{}
}

func (i *inmemory) Save(sum float64, category string) (model.Spend, error) {
	rec := model.Spend{
		ID:       i.lastIndex + 1,
		Price:    sum,
		Category: category,
		DateTime: time.Now(),
	}
	i.records = append(i.records, rec)

	return rec, nil
}

func (i *inmemory) GetByTimeSince(timeSince time.Time) ([]model.Spend, error) {
	var result []model.Spend
	for _, rec := range i.records {
		if timeSince.Before(rec.DateTime) {
			result = append(result, rec)
		}
	}

	return result, nil
}
