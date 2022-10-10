package command

import (
	currencyRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency"
	spendRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/command/currency"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/command/report"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/command/spend"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/command/start"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/command/unknown"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages"
)

func MakeChain(spendRepo spendRepo.Repository,
	currencyRepo currencyRepo.Repository,
	selectedCurrencyRepo selected_currency.Repository,
) messages.Command {
	unknownCmd := unknown.New()
	currencyCmd := currency.New(unknownCmd, selectedCurrencyRepo)
	spendCmd := spend.New(currencyCmd, spendRepo)
	reportCmd := report.New(spendCmd, spendRepo, selectedCurrencyRepo, currencyRepo)

	return start.New(reportCmd)
}
