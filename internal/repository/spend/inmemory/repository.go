package inmemory

import (
	"time"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
)

type inmemory struct {
	lastIndex int64
	records   []spend.Record
}

func New() *inmemory {
	return &inmemory{
		lastIndex: 0,
		records:   []spend.Record{},
	}
}

func (i *inmemory) Save(sum int64, category string) (spend.Record, error) {
	rec := spend.Record{
		ID:       i.lastIndex + 1,
		Sum:      sum,
		Category: category,
		DateTime: time.Now(),
	}
	i.records = append(i.records, rec)

	return rec, nil
}

func (i *inmemory) GetByTimeSince(timeSince time.Time) ([]spend.Record, error) {
	var result []spend.Record
	for _, rec := range i.records {
		if timeSince.Before(rec.DateTime) {
			result = append(result, rec)
		}
	}

	return result, nil
}
