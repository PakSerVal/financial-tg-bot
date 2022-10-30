package report

import (
	"context"
	"time"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/report"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/utils"
)

const (
	commandToday = "today"
	commandMonth = "month"
	commandYear  = "year"
)

type reportCommand struct {
	next          messages.Command
	reportService report.Service
}

func New(
	next messages.Command,
	reportService report.Service,
) messages.Command {
	return &reportCommand{
		next:          next,
		reportService: reportService,
	}
}

func (r *reportCommand) Process(ctx context.Context, in model.MessageIn) (*model.MessageOut, error) {
	now := time.Now()
	switch in.Command {
	case commandToday:
		return r.reportService.MakeReport(
			ctx,
			in.UserId,
			utils.BeginOfDay(now),
			"сегодня",
		)
	case commandMonth:
		return r.reportService.MakeReport(
			ctx,
			in.UserId,
			utils.BeginOfMonth(now),
			"в текущем месяце",
		)
	case commandYear:
		return r.reportService.MakeReport(
			ctx,
			in.UserId,
			utils.BeginOfYear(now),
			"в этом году",
		)
	}

	return r.next.Process(ctx, in)
}

func (r *reportCommand) Name() string {
	return "report"
}
