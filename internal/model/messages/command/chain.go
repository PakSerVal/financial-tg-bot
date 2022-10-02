package command

import (
	"time"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/report"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/spend"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/start"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/unknown"
	spend_repo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
)

type Repository interface {
	Save(sum int64, category string) (spend_repo.SpendRecord, error)
	GetByTimeSince(timeSince time.Time) ([]spend_repo.SpendRecord, error)
}

func MakeChain(repo Repository) messages.Command {
	return start.New(spend.New(report.New(unknown.New(), repo), repo))
}
