package command

import (
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/currency"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/report"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/spend"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/start"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/unknown"
	currencyRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/currency_rate"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency"
	spendRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
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
