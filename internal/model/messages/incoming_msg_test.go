package messages

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "gitlab.ozon.dev/paksergey94/telegram-bot/internal/mocks/messages"
)

func Test_CommandProcessError(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mocks.NewMockMessageSender(ctrl)
	chain := mocks.NewMockCommand(ctrl)
	model := New(sender, chain)

	chain.EXPECT().Process("/start").Return("", errors.New("some error"))

	err := model.IncomingMessage(Message{
		Text:   "/start",
		UserID: 123,
	})

	assert.Error(t, err)
}

func Test_CommandProcessSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)

	sender := mocks.NewMockMessageSender(ctrl)
	sender.EXPECT().SendMessage("привет", int64(123))
	chain := mocks.NewMockCommand(ctrl)
	chain.EXPECT().Process("some text").Return("привет", nil)
	model := New(sender, chain)

	err := model.IncomingMessage(Message{
		Text:   "some text",
		UserID: 123,
	})

	assert.NoError(t, err)
}
