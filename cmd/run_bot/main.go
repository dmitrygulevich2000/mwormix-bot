package main

import (
	"errors"
	"log"
	"os"

	"github.com/dmitrygulevich2000/mwormix-bot/internal/bot"
)

func RunApp() error {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		return errors.New("env TELEGRAM_TOKEN not set")
	}

	tgbot, err := bot.NewTelegram(token)
	if err != nil {
		return err
	}

	return tgbot.Run()
}

func main() {
	err := RunApp()
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Shutted down")
	}

}
