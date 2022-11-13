package report

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/report/queue_message"
)

const (
	commandToday = "today"
	commandMonth = "month"
	commandYear  = "year"
)

type reportCommand struct {
	next         messages.Command
	reportSender queue_message.Sender
}

func New(
	next messages.Command,
	reportSender queue_message.Sender,
) messages.Command {
	return &reportCommand{
		next:         next,
		reportSender: reportSender,
	}
}

func (r *reportCommand) Process(ctx context.Context, in model.MessageIn) (*model.MessageOut, error) {
	if in.Command != commandToday && in.Command != commandMonth && in.Command != commandYear {
		return r.next.Process(ctx, in)
	}

	msg := model.ReportMsg{
		UserId: in.UserId,
		Period: in.Command,
	}

	err := r.reportSender.Send(msg)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return nil, nil
}

func (r *reportCommand) Name() string {
	return "report"
}
