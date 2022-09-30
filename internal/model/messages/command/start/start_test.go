package start

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mocks "gitlab.ozon.dev/paksergey94/telegram-bot/internal/mocks/messages"
)

func Test_NotSupported(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mocks.NewMockCommand(ctrl)

	command := New(next)

	next.EXPECT().Process("not supported text").Return("привет", nil)

	res, err := command.Process("not supported text")

	assert.NoError(t, err)
	assert.Equal(t, "привет", res)
}

func Test_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mocks.NewMockCommand(ctrl)

	command := New(next)

	res, err := command.Process("/start")

	assert.NoError(t, err)
	assert.Equal(
		t,
		"Бот для учета финансов\n\nДобавить трату: 350 продукты\n\nПолучить отчет: \n- за сегодня: /today\n- за месяц: /month\n- за год: /year\n",
		res)
}
