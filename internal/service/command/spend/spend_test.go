package spend

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/database"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
	mockBudget "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/budget/mocks"
	mock_cache "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend/cache/mocks"
	mockSpend "gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend/mocks"
	mockMessages "gitlab.ozon.dev/paksergey94/telegram-bot/internal/service/messages/mocks"
)

func TestSpendCommand_Process(t *testing.T) {
	db, mock, _ := sqlmock.New()

	ctrl := gomock.NewController(t)
	next := mockMessages.NewMockCommand(ctrl)
	spendRepo := mockSpend.NewMockRepository(ctrl)
	budgetRepo := mockBudget.NewMockRepository(ctrl)
	sqlManager := database.NewSqlManager(db)
	cache := mock_cache.NewMockSpendRepo(ctrl)

	command := New(next, spendRepo, budgetRepo, sqlManager, cache)

	t.Run("not supported", func(t *testing.T) {
		next.EXPECT().Process(context.TODO(), model.MessageIn{Command: "not supported text", UserId: 123}).Return(&model.MessageOut{Text: "привет"}, nil)

		res, err := command.Process(context.TODO(), model.MessageIn{Command: "not supported text", UserId: 123})

		assert.NoError(t, err)
		assert.Equal(t, &model.MessageOut{Text: "привет"}, res)
	})

	t.Run("spendRepo save error", func(t *testing.T) {
		mock.ExpectBegin()
		spendRepo.EXPECT().SaveTx(gomock.Any(), context.TODO(), int64(12300), "такси", int64(123)).Return(errors.New("some error"))
		mock.ExpectRollback()

		_, err := command.Process(context.TODO(), model.MessageIn{Command: "123 такси", UserId: 123})

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("spendRepo get by time since error", func(t *testing.T) {
		mock.ExpectBegin()
		spendRepo.EXPECT().SaveTx(gomock.Any(), context.TODO(), int64(12300), "такси", int64(123)).Return(nil)
		spendRepo.EXPECT().GetByTimeSinceTx(gomock.Any(), context.TODO(), int64(123), gomock.Any()).Return(nil, errors.New("some error"))
		mock.ExpectRollback()

		_, err := command.Process(context.TODO(), model.MessageIn{Command: "123 такси", UserId: 123})

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("budget exceeded", func(t *testing.T) {
		mock.ExpectBegin()
		spendRepo.EXPECT().SaveTx(gomock.Any(), context.TODO(), int64(100), "такси", int64(123)).Return(nil)
		spendRepo.EXPECT().GetByTimeSinceTx(gomock.Any(), context.TODO(), int64(123), gomock.Any()).Return([]model.Spend{
			{
				Id:       1,
				Price:    3000000,
				Category: "cat1",
				UserId:   123,
			},
			{
				Id:       2,
				Price:    1,
				Category: "cat2",
				UserId:   123,
			},
		}, nil)
		budgetRepo.EXPECT().GetBudgetTx(gomock.Any(), context.TODO(), int64(123)).Return(&model.Budget{
			Value: 3000000,
		}, nil)
		mock.ExpectRollback()

		cache.EXPECT().DeleteForUser(context.TODO(), int64(123)).Return(nil)

		res, err := command.Process(context.TODO(), model.MessageIn{Command: "1 такси", UserId: 123})

		assert.NoError(t, err)
		assert.Equal(t, &model.MessageOut{Text: "Трата не была добавлена, так как превышен лимит за текущий месяц в 30000.00 руб"}, res)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("save to cache error", func(t *testing.T) {
		mock.ExpectBegin()
		spendRepo.EXPECT().SaveTx(gomock.Any(), context.TODO(), int64(12300), "такси", int64(123)).Return(nil)
		spendRepo.EXPECT().GetByTimeSinceTx(gomock.Any(), context.TODO(), int64(123), gomock.Any()).Return([]model.Spend{
			{
				Id:       1,
				Price:    30000,
				Category: "cat1",
				UserId:   123,
			},
			{
				Id:       2,
				Price:    1,
				Category: "cat2",
				UserId:   123,
			},
		}, nil)
		budgetRepo.EXPECT().GetBudgetTx(gomock.Any(), context.TODO(), int64(123)).Return(nil, nil)
		mock.ExpectCommit()

		cache.EXPECT().DeleteForUser(context.TODO(), int64(123)).Return(errors.New("some error"))

		res, err := command.Process(context.TODO(), model.MessageIn{Command: "123 такси", UserId: 123})

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		spendRepo.EXPECT().SaveTx(gomock.Any(), context.TODO(), int64(12300), "такси", int64(123)).Return(nil)
		spendRepo.EXPECT().GetByTimeSinceTx(gomock.Any(), context.TODO(), int64(123), gomock.Any()).Return([]model.Spend{
			{
				Id:       1,
				Price:    30000,
				Category: "cat1",
				UserId:   123,
			},
			{
				Id:       2,
				Price:    1,
				Category: "cat2",
				UserId:   123,
			},
		}, nil)
		budgetRepo.EXPECT().GetBudgetTx(gomock.Any(), context.TODO(), int64(123)).Return(nil, nil)
		mock.ExpectCommit()

		cache.EXPECT().DeleteForUser(context.TODO(), int64(123)).Return(nil)

		res, err := command.Process(context.TODO(), model.MessageIn{Command: "123 такси", UserId: 123})

		assert.NoError(t, err)
		assert.Equal(t, &model.MessageOut{Text: "Добавлена трата: такси 123.00 руб."}, res)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
