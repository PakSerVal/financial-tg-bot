package spend

import "time"

type Repository interface {
	Save(sum float64, category string) (SpendRecord, error)
	GetByTimeSince(timeSince time.Time) ([]SpendRecord, error)
}
