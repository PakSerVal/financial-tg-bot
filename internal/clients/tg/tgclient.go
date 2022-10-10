package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
)

const (
	updateOffset  = 0
	updateTimeout = 60
)

type Client interface {
	SendMessage(msgOut model.MessageOut, userID int64) error
	GetUpdatesChan() tgbotapi.UpdatesChannel
}

type TokenGetter interface {
	Token() string
}

type client struct {
	client *tgbotapi.BotAPI
}

func New(tokenGetter TokenGetter) (Client, error) {
	c, err := tgbotapi.NewBotAPI(tokenGetter.Token())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new bot")
	}

	return &client{
		client: c,
	}, nil
}

func (c *client) SendMessage(msg model.MessageOut, userID int64) error {
	_, err := c.client.Send(makeTgMessage(msg, userID))
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}
	return nil
}

func (c *client) GetUpdatesChan() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(updateOffset)
	u.Timeout = updateTimeout

	return c.client.GetUpdatesChan(u)
}

func makeTgMessage(msg model.MessageOut, userID int64) tgbotapi.MessageConfig {
	tgMessage := tgbotapi.NewMessage(userID, msg.Text)
	if msg.KeyBoard != nil {
		tgRows := make([][]tgbotapi.KeyboardButton, 0, len(msg.KeyBoard.Rows))
		for _, row := range msg.KeyBoard.Rows {
			tgButtons := make([]tgbotapi.KeyboardButton, 0, len(row.Buttons))

			for _, button := range row.Buttons {
				tgButtons = append(tgButtons, tgbotapi.NewKeyboardButton(button.Text))
			}

			tgRows = append(tgRows, tgButtons)
		}

		if msg.KeyBoard.OneTime {
			tgMessage.ReplyMarkup = tgbotapi.NewOneTimeReplyKeyboard(tgRows...)
		} else {
			tgMessage.ReplyMarkup = tgbotapi.NewReplyKeyboard(tgRows...)
		}
	}

	return tgMessage
}
