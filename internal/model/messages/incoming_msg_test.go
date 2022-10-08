package messages

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "gitlab.ozon.dev/paksergey94/telegram-bot/internal/mocks/messages"
)

func TestModel_ProcessMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	chain := mocks.NewMockCommand(ctrl)
	model := New(sender, chain)

	t.Run("chain error", func(t *testing.T) {
		chain.EXPECT().Process("/start").Return("", errors.New("some error"))
		err := model.processMessage(Message{Text: "/start"})
		assert.Error(t, err)
	})

	t.Run("sender error", func(t *testing.T) {
		chain.EXPECT().Process("/start").Return("привет", nil)
		sender.EXPECT().SendMessage("привет", int64(1)).Return(errors.New("some error"))

		err := model.processMessage(Message{Text: "/start", UserID: 1})
		assert.Error(t, err)
	})
}
