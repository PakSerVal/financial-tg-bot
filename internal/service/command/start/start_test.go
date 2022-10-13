package start

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	mockMessages "gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages/mocks"
)

func TestStartCommand_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mockMessages.NewMockCommand(ctrl)

	next.EXPECT().Process(model.MessageIn{Text: "not supported text"}).Return(&model.MessageOut{Text: "привет"}, nil)

	command := New(next)

	t.Run("not supported text", func(t *testing.T) {
		res, err := command.Process(model.MessageIn{Text: "not supported text"})

		assert.NoError(t, err)
		assert.Equal(t, &model.MessageOut{Text: "привет"}, res)
	})

	t.Run("success", func(t *testing.T) {
		res, err := command.Process(model.MessageIn{Text: "start"})

		assert.NoError(t, err)
		assert.Equal(
			t,
			&model.MessageOut{
				Text:     "Бот для учета финансов\n\nДобавить трату: 350 продукты\n\nИзменить валюту: /currency\n\nПолучить отчет: \n- за сегодня: /today\n- за месяц: /month\n- за год: /year\n",
				KeyBoard: nil,
			},
			res)
	})
}
