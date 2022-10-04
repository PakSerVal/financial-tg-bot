package start

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "gitlab.ozon.dev/paksergey94/telegram-bot/internal/mocks/messages"
)

func TestStartCommand_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mocks.NewMockCommand(ctrl)

	next.EXPECT().Process("not supported text").Return("привет", nil)

	command := New(next)

	t.Run("not supported text", func(t *testing.T) {
		res, err := command.Process("not supported text")

		assert.NoError(t, err)
		assert.Equal(t, "привет", res)
	})

	t.Run("success", func(t *testing.T) {
		res, err := command.Process("start")

		assert.NoError(t, err)
		assert.Equal(
			t,
			"Бот для учета финансов\n\nДобавить трату: 350 продукты\n\nПолучить отчет: \n- за сегодня: /today\n- за месяц: /month\n- за год: /year\n",
			res)
	})
}
