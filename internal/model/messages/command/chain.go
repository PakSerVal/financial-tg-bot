package command

import (
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/report"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/spend"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/start"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/unknown"
)

func MakeChain(repo spend.Repository) messages.Command {
	return start.New(spend.New(report.New(unknown.New(), repo), repo))
}
