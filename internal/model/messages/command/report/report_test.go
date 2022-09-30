package report

import (
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	mocks "gitlab.ozon.dev/paksergey94/telegram-bot/internal/mocks/messages"
	mock_report "gitlab.ozon.dev/paksergey94/telegram-bot/internal/mocks/messages/command/report"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
)

func Test_NotSupported(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mocks.NewMockCommand(ctrl)
	repo := mock_report.NewMockRepository(ctrl)

	command := New(next, repo)

	next.EXPECT().Process("not supported text").Return("привет", nil)

	res, err := command.Process("not supported text")

	assert.NoError(t, err)
	assert.Equal(t, "привет", res)
}

func Test_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mocks.NewMockCommand(ctrl)
	repo := mock_report.NewMockRepository(ctrl)
	repo.EXPECT().GetByTimeSince(gomock.Any()).Return([]spend.Record{}, errors.New("some error"))

	command := New(next, repo)

	_, err := command.Process("/today")

	assert.Error(t, err)
}

func Test_NoRecords(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mocks.NewMockCommand(ctrl)
	repo := mock_report.NewMockRepository(ctrl)
	repo.EXPECT().GetByTimeSince(gomock.Any()).Return([]spend.Record{}, nil)

	command := New(next, repo)

	res, err := command.Process("/today")

	assert.NoError(t, err)
	assert.Equal(t, "Расходов сегодня нет", res)
}

func Test_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mocks.NewMockCommand(ctrl)
	repo := mock_report.NewMockRepository(ctrl)
	repo.EXPECT().GetByTimeSince(gomock.Any()).Times(3).Return([]spend.Record{
		{
			ID:       1,
			Sum:      100,
			Category: "Такси",
		},
		{
			ID:       2,
			Sum:      400,
			Category: "Такси",
		},
		{
			ID:       3,
			Sum:      200,
			Category: "Такси",
		},
		{
			ID:       4,
			Sum:      200,
			Category: "Продукты",
		},
		{
			ID:       4,
			Sum:      900,
			Category: "Продукты",
		},
		{
			ID:       5,
			Sum:      2000,
			Category: "Инвестиции",
		},
	}, nil)

	cases := []struct {
		command string
		wanted  string
	}{
		{
			command: "/today",
			wanted:  "Расходы сегодня:\nИнвестиции - 2000 руб.\nПродукты - 1100 руб.\nТакси - 700 руб.",
		},
		{
			command: "/month",
			wanted:  "Расходы в текущем месяце:\nИнвестиции - 2000 руб.\nПродукты - 1100 руб.\nТакси - 700 руб.",
		},
		{
			command: "/year",
			wanted:  "Расходы в этом году:\nИнвестиции - 2000 руб.\nПродукты - 1100 руб.\nТакси - 700 руб.",
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			command := New(next, repo)
			res, err := command.Process(c.command)

			assert.NoError(t, err)
			assert.Equal(t, c.wanted, res)
		})
	}
}
