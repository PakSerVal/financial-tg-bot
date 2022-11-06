package command

import (
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/database"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/budget"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/selected_currency"
	spendRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend/cache"
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
	spendCache cache.SpendRepo,
) messages.Command {
	unknownCmd := WithObserve(unknown.New())
	currencyCmd := WithObserve(currency.New(unknownCmd, selectedCurrencyRepo))
	spendCmd := WithObserve(spend.New(currencyCmd, spendRepo, budgetRepo, sqlManager, spendCache))
	reportCmd := WithObserve(report.New(spendCmd, reportService))
	budgetCmd := WithObserve(budgetCommand.New(reportCmd, budgetRepo))

	return start.New(budgetCmd)
}
