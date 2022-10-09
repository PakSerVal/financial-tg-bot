package inmemory

import (
	"time"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
)

type inmemory struct {
	lastIndex int64
	records   []spend.SpendRecord
}

func New() spend.Repository {
	return &inmemory{}
}

func (i *inmemory) Save(sum float64, category string) (spend.SpendRecord, error) {
	rec := spend.SpendRecord{
		ID:       i.lastIndex + 1,
		Price:    sum,
		Category: category,
		DateTime: time.Now(),
	}
	i.records = append(i.records, rec)

	return rec, nil
}

func (i *inmemory) GetByTimeSince(timeSince time.Time) ([]spend.SpendRecord, error) {
	var result []spend.SpendRecord
	for _, rec := range i.records {
		if timeSince.Before(rec.DateTime) {
			result = append(result, rec)
		}
	}

	return result, nil
}
