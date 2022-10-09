package messages

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "gitlab.ozon.dev/paksergey94/telegram-bot/internal/mocks/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command/dto"
)

func TestModel_ProcessMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	chain := mocks.NewMockCommand(ctrl)
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
