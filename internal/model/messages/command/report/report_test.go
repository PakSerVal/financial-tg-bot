package report

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	mocks "gitlab.ozon.dev/paksergey94/telegram-bot/internal/mocks/messages"
	mock_report "gitlab.ozon.dev/paksergey94/telegram-bot/internal/mocks/messages/command/report"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
)

func TestReportCommand_ProcessFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mocks.NewMockCommand(ctrl)
	repo := mock_report.NewMockRepository(ctrl)
	command := New(next, repo)

	gomock.InOrder(
		next.EXPECT().Process("not supported text").Return("привет", nil),
		repo.EXPECT().GetByTimeSince(gomock.Any()).Return([]spend.SpendRecord{}, errors.New("some error")).Times(1),
		repo.EXPECT().GetByTimeSince(gomock.Any()).Return([]spend.SpendRecord{}, nil).Times(1),
		repo.EXPECT().GetByTimeSince(gomock.Any()).Times(3).Return([]spend.SpendRecord{
			{
				ID:       1,
				Price:    100,
				Category: "Такси",
			},
			{
				ID:       2,
				Price:    400,
				Category: "Такси",
			},
			{
				ID:       3,
				Price:    200,
				Category: "Такси",
			},
			{
				ID:       4,
				Price:    200,
				Category: "Продукты",
			},
			{
				ID:       4,
				Price:    900,
				Category: "Продукты",
			},
			{
				ID:       5,
				Price:    2000,
				Category: "Инвестиции",
			},
		}, nil),
	)

	t.Run("not supported", func(t *testing.T) {
		res, err := command.Process("not supported text")

		assert.NoError(t, err)
		assert.Equal(t, "привет", res)
	})

	t.Run("repo error", func(t *testing.T) {
		_, err := command.Process("/today")

		assert.Error(t, err)
	})

	t.Run("no records", func(t *testing.T) {
		res, err := command.Process("/today")

		assert.NoError(t, err)
		assert.Equal(t, "Расходов сегодня нет", res)
	})

	sucessCases := []struct {
		name    string
		command string
		wanted  string
	}{
		{
			name:    "today",
			command: "/today",
			wanted:  "Расходы сегодня:\nИнвестиции - 2000 руб.\nПродукты - 1100 руб.\nТакси - 700 руб.",
		},
		{
			name:    "month",
			command: "/month",
			wanted:  "Расходы в текущем месяце:\nИнвестиции - 2000 руб.\nПродукты - 1100 руб.\nТакси - 700 руб.",
		},
		{
			name:    "year",
			command: "/year",
			wanted:  "Расходы в этом году:\nИнвестиции - 2000 руб.\nПродукты - 1100 руб.\nТакси - 700 руб.",
		},
	}

	for _, c := range sucessCases {
		t.Run(c.name, func(t *testing.T) {
			res, err := command.Process(c.command)

			assert.NoError(t, err)
			assert.Equal(t, c.wanted, res)
		})
	}
}
