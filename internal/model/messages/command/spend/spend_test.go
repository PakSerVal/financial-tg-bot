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

func Test_NotSupported(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mocks.NewMockCommand(ctrl)
	repo := mock_spend.NewMockRepository(ctrl)

	command := New(next, repo)

	next.EXPECT().Process("not supported text").Return("привет", nil)

	res, err := command.Process("not supported text")

	assert.NoError(t, err)
	assert.Equal(t, "привет", res)
}

func Test_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mocks.NewMockCommand(ctrl)
	repo := mock_spend.NewMockRepository(ctrl)

	command := New(next, repo)

	repo.EXPECT().Save(int64(123), "такси").Return(spend.Record{}, errors.New("some error"))

	_, err := command.Process("123 такси")

	assert.Error(t, err)
}

func Test_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	next := mocks.NewMockCommand(ctrl)
	repo := mock_spend.NewMockRepository(ctrl)
	repo.EXPECT().Save(int64(123), "такси").Return(spend.Record{
		ID:       1,
		Sum:      123,
		Category: "Такси",
	}, nil)

	command := New(next, repo)

	res, err := command.Process("123 такси")

	assert.NoError(t, err)
	assert.Equal(t, "Добавлена трата: Такси 123 руб.", res)
}
