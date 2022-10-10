package messages

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/dto"
	mockMessages "gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/mocks"
)

func TestModel_ProcessMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mockMessages.NewMockMessageSender(ctrl)
	chain := mockMessages.NewMockCommand(ctrl)
	model := New(sender, chain)

	t.Run("chain error", func(t *testing.T) {
		chain.EXPECT().Process(dto.MessageIn{Text: "/start", UserId: 1}).Return(dto.MessageOut{}, errors.New("some error"))
		err := model.processMessage("/start", 1)
		assert.Error(t, err)
	})

	t.Run("sender error", func(t *testing.T) {
		chain.EXPECT().Process(dto.MessageIn{Text: "/start", UserId: 1}).Return(dto.MessageOut{Text: "привет"}, nil)
		sender.EXPECT().SendMessage(dto.MessageOut{Text: "привет"}, int64(1)).Return(errors.New("some error"))

		err := model.processMessage("/start", int64(1))
		assert.Error(t, err)
	})
}
