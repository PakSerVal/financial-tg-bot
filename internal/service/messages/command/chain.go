package command

import (
	currencyRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency"
	spendRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages/command/currency"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages/command/report"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages/command/spend"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages/command/start"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages/command/unknown"
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
