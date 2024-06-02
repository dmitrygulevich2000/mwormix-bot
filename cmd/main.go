package main

import (
	"log"
	"os"

	"github.com/dmitrygulevich2000/mwormix-bot/internal/bot"
	"github.com/dmitrygulevich2000/mwormix-bot/internal/collector"
	"github.com/dmitrygulevich2000/mwormix-bot/internal/storage"
	"github.com/dmitrygulevich2000/mwormix-bot/internal/utils"
)

func main() {
	// req, err := http.NewRequest("GET", "https://vk.com/wormix_club", nil)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	// resp, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// f, _ := os.Create("testdata.html")
	// io.Copy(f, resp.Body)

	/////////////////////////////////////////////

	server := utils.RunFileServer("testdata.html")

	searches, err := storage.NewSearches("tmpdb")
	defer os.Remove("tmpdb")
	if err != nil {
		log.Fatalln(err)
	}
	storage := storage.New(searches)

	collector := collector.New(storage, server.URL)
	bonuses, err := collector.CollectNewBonuses()
	if err != nil {
		log.Fatalln(err)
	}

	/////////////////////////////////////////////

	token := os.Getenv("MWORMIX_TELEGRAM_TOKEN")
	if token == "" {
		log.Fatalln("env MWORMIX_TELEGRAM_TOKEN not set")
	}

	tgbot, err := bot.NewTelegram(token)
	if err != nil {
		log.Fatalln(err)
	}

	tgbot.Broadcast(&bonuses[0])

	// err = tgbot.Run()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// println("Shutted down")
}
