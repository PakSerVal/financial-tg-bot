package spend

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	mocks "gitlab.ozon.dev/paksergey94/telegram-bot/internal/mocks/messages"
	mock_spend "gitlab.ozon.dev/paksergey94/telegram-bot/internal/mocks/messages/command/spend"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
)

func TestSpendCommand_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mocks.NewMockCommand(ctrl)
	repo := mock_spend.NewMockRepository(ctrl)

	command := New(next, repo)

	gomock.InOrder(
		next.EXPECT().Process("not supported text").Return("привет", nil),
		repo.EXPECT().Save(int64(123), "такси").Return(spend.SpendRecord{}, errors.New("some error")),
		repo.EXPECT().Save(int64(123), "такси").Return(spend.SpendRecord{
			ID:       1,
			Price:    123,
			Category: "Такси",
		}, nil),
	)

	t.Run("not supported", func(t *testing.T) {
		res, err := command.Process("not supported text")

		assert.NoError(t, err)
		assert.Equal(t, "привет", res)
	})

	t.Run("repo error", func(t *testing.T) {
		_, err := command.Process("123 такси")

		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		res, err := command.Process("123 такси")

		assert.NoError(t, err)
		assert.Equal(t, "Добавлена трата: Такси 123 руб.", res)
	})
}
