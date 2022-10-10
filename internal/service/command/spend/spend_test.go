package spend

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	mockSpend "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend/mocks"
	mockMessages "gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages/mocks"
)

func TestSpendCommand_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mockMessages.NewMockCommand(ctrl)
	sendRepo := mockSpend.NewMockRepository(ctrl)

	command := New(next, sendRepo)

	gomock.InOrder(
		next.EXPECT().Process(model.MessageIn{Text: "not supported text"}).Return(model.MessageOut{Text: "привет"}, nil),
		sendRepo.EXPECT().Save(float64(123), "такси").Return(model.Spend{}, errors.New("some error")),
		sendRepo.EXPECT().Save(float64(123), "такси").Return(model.Spend{
			ID:       1,
			Price:    123,
			Category: "Такси",
		}, nil),
	)

	t.Run("not supported", func(t *testing.T) {
		res, err := command.Process(model.MessageIn{Text: "not supported text"})

		assert.NoError(t, err)
		assert.Equal(t, model.MessageOut{Text: "привет"}, res)
	})

	t.Run("repo error", func(t *testing.T) {
		_, err := command.Process(model.MessageIn{Text: "123 такси"})

		assert.Error(t, err)
	})

	t.Run("success", func(t *testing.T) {
		res, err := command.Process(model.MessageIn{Text: "123 такси"})

		assert.NoError(t, err)
		assert.Equal(t, model.MessageOut{Text: "Добавлена трата: Такси 123.00 руб."}, res)
	})
}
