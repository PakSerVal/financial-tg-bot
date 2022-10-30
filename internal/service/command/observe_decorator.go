package command

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/metrics"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages"
)

type observeDecorator struct {
	command messages.Command
}

func WithObserve(command messages.Command) messages.Command {
	return &observeDecorator{
		command: command,
	}
}

func (t *observeDecorator) Process(ctx context.Context, in model.MessageIn) (*model.MessageOut, error) {
	start := time.Now()

	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"message process",
	)
	defer span.Finish()

	span.SetTag("command", t.command.Name())

	defer metrics.MessageProcessedTime(time.Since(start).Seconds(), t.command.Name())

	return t.command.Process(ctx, in)
}

func (t *observeDecorator) Name() string {
	return "observe_decorator"
}
