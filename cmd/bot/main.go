package main

import (
	"log"

	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/config"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model/messages/command"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend/inmemory"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	tgClient, err := tg.New(cfg)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}
	repo := inmemory.New()
	msgModel := messages.New(tgClient, command.MakeChain(repo))

	tgClient.ListenUpdates(msgModel)
}
