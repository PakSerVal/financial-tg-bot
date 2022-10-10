package messages

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_tg "gitlab.ozon.dev/paksergey94/telegram-bot/internal/clients/tg/mocks"
	model2 "gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	mockMessages "gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages/mocks"
)

func TestModel_ProcessMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mock_tg.NewMockClient(ctrl)
	chain := mockMessages.NewMockCommand(ctrl)
	model := New(sender, chain)

	t.Run("chain error", func(t *testing.T) {
		chain.EXPECT().Process(model2.MessageIn{Text: "/start", UserId: 1}).Return(model2.MessageOut{}, errors.New("some error"))
		err := model.processMessage("/start", 1)
		assert.Error(t, err)
	})

	t.Run("sender error", func(t *testing.T) {
		chain.EXPECT().Process(model2.MessageIn{Text: "/start", UserId: 1}).Return(model2.MessageOut{Text: "привет"}, nil)
		sender.EXPECT().SendMessage(model2.MessageOut{Text: "привет"}, int64(1)).Return(errors.New("some error"))

		err := model.processMessage("/start", int64(1))
		assert.Error(t, err)
	})
}
