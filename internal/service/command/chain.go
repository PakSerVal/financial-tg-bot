package command

import (
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency"
	spendRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/command/currency"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/command/report"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/command/spend"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/command/start"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/command/unknown"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages"
	reportService "gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/report"
)

func MakeChain(spendRepo spendRepo.Repository,
	selectedCurrencyRepo selected_currency.Repository,
	reportService reportService.Service,
) messages.Command {
	unknownCmd := unknown.New()
	currencyCmd := currency.New(unknownCmd, selectedCurrencyRepo)
	spendCmd := spend.New(currencyCmd, spendRepo)
	reportCmd := report.New(spendCmd, reportService)

	return start.New(reportCmd)
}
