package command

import (
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/database"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/budget"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency"
	spendRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
	budgetCommand "gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/command/budget"
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
	budgetRepo budget.Repository,
	sqlManager database.SqlManager,
) messages.Command {
	unknownCmd := unknown.New()
	currencyCmd := currency.New(unknownCmd, selectedCurrencyRepo)
	spendCmd := spend.New(currencyCmd, spendRepo, budgetRepo, sqlManager)
	reportCmd := report.New(spendCmd, reportService)
	budgetCmd := budgetCommand.New(reportCmd, budgetRepo)

	return start.New(budgetCmd)
}
