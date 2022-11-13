package report

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	mockMessages "gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages/mocks"
	mock_queue_message "gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/report/queue_message/mocks"
)

func TestReportCommand_ProcessFailed(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	next := mockMessages.NewMockCommand(ctrl)
	sender := mock_queue_message.NewMockSender(ctrl)

	command := New(next, sender)

	t.Run("not supported", func(t *testing.T) {
		next.EXPECT().Process(ctx, model.MessageIn{Command: "not supported"}).Return(&model.MessageOut{Text: "test"}, nil)
		res, err := command.Process(ctx, model.MessageIn{Command: "not supported"})

		assert.NoError(t, err)
		assert.Equal(t, &model.MessageOut{Text: "test"}, res)
	})

	t.Run("sender error", func(t *testing.T) {
		sender.EXPECT().Send(gomock.Any()).Return(errors.New("some error"))
		res, err := command.Process(ctx, model.MessageIn{Command: "today", UserId: 1})

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("sender success", func(t *testing.T) {
		sender.EXPECT().Send(model.ReportMsg{
			UserId: 1,
			Period: "today",
		}).Return(nil)

		res, err := command.Process(ctx, model.MessageIn{Command: "today", UserId: 1})

		assert.NoError(t, err)
		assert.Nil(t, res)
	})
}
