package report

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	currencyRepo "gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	mockMessages "gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages/mocks"
	mock_report "gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/report/mocks"
)

func TestReportCommand_ProcessFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mockMessages.NewMockCommand(ctrl)
	service := mock_report.NewMockService(ctrl)

	command := New(next, service)

	t.Run("not supported", func(t *testing.T) {
		next.EXPECT().Process(currencyRepo.MessageIn{Text: "not supported"}).Return(&currencyRepo.MessageOut{Text: "test"}, nil)
		res, err := command.Process(currencyRepo.MessageIn{Text: "not supported"})

		assert.NoError(t, err)
		assert.Equal(t, &currencyRepo.MessageOut{Text: "test"}, res)
	})

	t.Run("today report", func(t *testing.T) {
		service.EXPECT().MakeReport(int64(1), gomock.Any(), "сегодня").Return(&currencyRepo.MessageOut{Text: "test report"}, nil)
		res, err := command.Process(currencyRepo.MessageIn{Text: "today", UserId: 1})

		assert.NoError(t, err)
		assert.Equal(t, &currencyRepo.MessageOut{Text: "test report"}, res)
	})

	t.Run("month report", func(t *testing.T) {
		service.EXPECT().MakeReport(int64(1), gomock.Any(), "в текущем месяце").Return(&currencyRepo.MessageOut{Text: "test report"}, nil)
		res, err := command.Process(currencyRepo.MessageIn{Text: "month", UserId: 1})

		assert.NoError(t, err)
		assert.Equal(t, &currencyRepo.MessageOut{Text: "test report"}, res)
	})

	t.Run("year report", func(t *testing.T) {
		service.EXPECT().MakeReport(int64(1), gomock.Any(), "в этом году").Return(&currencyRepo.MessageOut{Text: "test report"}, nil)
		res, err := command.Process(currencyRepo.MessageIn{Text: "year", UserId: 1})

		assert.NoError(t, err)
		assert.Equal(t, &currencyRepo.MessageOut{Text: "test report"}, res)
	})
}
