package bot

import (
	"fmt"
	"slices"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/dmitrygulevich2000/mwormix-bot/internal/collector"
)

type Telegram struct {
	botAPI *tgbotapi.BotAPI
}

func NewTelegram(token string) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Telegram{
		botAPI: bot,
	}, nil
}

const chatID = 637215348

func (t *Telegram) BroadcastAll(bonuses []collector.Bonus) error {
	slices.SortFunc[[]collector.Bonus, collector.Bonus](bonuses, func(lhs, rhs collector.Bonus) int {
		return lhs.PublishedAt.Compare(rhs.PublishedAt)
	})
	var lastError error
	for _, bonus := range bonuses {
		err := t.Broadcast(&bonus)
		if err != nil {
			lastError = err
		}
	}
	return lastError
}

func (t *Telegram) Broadcast(bonus *collector.Bonus) error {
	text := strings.Builder{}
	if bonus.Code == "" || bonus.ValidUntil == "" {
		text.WriteString("Неизвестный пост с промокодом:\n")
		text.WriteString(bonus.Link)
	} else {
		text.WriteString("Новый бонус в группе игры Вормикс от ")
		text.WriteString(bonus.PublishedAt.Format("02.01.2006 15:04"))
		text.WriteString("!\n")
		text.WriteString(bonus.Link)
		text.WriteString("\nПромокод: ")
		text.WriteString(bonus.Code)
		text.WriteString("\nБонус действителен ")
		text.WriteString(bonus.ValidUntil)
	}

	msg := tgbotapi.NewMessage(chatID, text.String())
	msg.DisableWebPagePreview = true

	_, err := t.botAPI.Send(msg)
	return err
}

func (t *Telegram) Run() error {
	t.botAPI.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "start":
			msg.Text = `Уведомляю о бонусах в группе игры Вормикс.
			
Поддерживаемые команды:
- /start
- /chat_id
`
		case "chat_id":
			msg.Text = fmt.Sprintf("this chat id: %d", update.Message.Chat.ID)
		default:
			msg.Text = "unknown command"
		}

		if _, err := t.botAPI.Send(msg); err != nil {
			return err
		}
	}
	return nil
}
