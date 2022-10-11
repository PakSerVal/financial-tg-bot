package report

import (
	"time"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/report"
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

func (r *reportCommand) Process(in model.MessageIn) (*model.MessageOut, error) {
	now := time.Now()
	switch in.Text {
	case commandToday:
		return r.reportService.MakeReport(
			in.UserId,
			time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
			"сегодня",
		)
	case commandMonth:
		return r.reportService.MakeReport(
			in.UserId,
			time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()),
			"в текущем месяце",
		)
	case commandYear:
		return r.reportService.MakeReport(
			in.UserId,
			time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location()),
			"в этом году",
		)
	}

	return r.next.Process(in)
}
