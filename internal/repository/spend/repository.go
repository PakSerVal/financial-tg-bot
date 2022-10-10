package spend

import (
	"time"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
)

type Repository interface {
	Save(sum float64, category string) (model.Spend, error)
	GetByTimeSince(timeSince time.Time) ([]model.Spend, error)
}
