package main

import (
	"errors"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"

	"github.com/dmitrygulevich2000/mwormix-bot/internal/bot"
	"github.com/dmitrygulevich2000/mwormix-bot/internal/collector"
	"github.com/dmitrygulevich2000/mwormix-bot/internal/storage"
)

func RunApp() error {
	var storageConf storage.Config
	err := envconfig.Process("mwormix", &storageConf)
	if err != nil {
		return err
	}

	storage, err := storage.NewFromConfig(storageConf)
	if err != nil {
		return err
	}
	bc := collector.New(storage, "https://vk.com/wormix_club")

	token := os.Getenv("MWORMIX_TELEGRAM_TOKEN")
	if token == "" {
		return errors.New("env MWORMIX_TELEGRAM_TOKEN not set")

	}
	tgbot, err := bot.NewTelegram(token)
	if err != nil {
		return err
	}

	bonuses, err := bc.CollectNewBonuses()
	if err != nil {
		return err
	}
	if len(bonuses) == 0 {
		log.Println("no new bonuses found")
	}

	return tgbot.BroadcastAll(bonuses)
}

func main() {
	err := RunApp()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
