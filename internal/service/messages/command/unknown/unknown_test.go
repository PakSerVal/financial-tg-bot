package unknown

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
)

func TestUnknownCommand_Process(t *testing.T) {
	command := New()

	res, err := command.Process(model.MessageIn{Text: "unknown command"})
	assert.NoError(t, err)
	assert.Equal(t, model.MessageOut{Text: "не знаю такую команду"}, res)
}
