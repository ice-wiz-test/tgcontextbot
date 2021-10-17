package handling

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func FindSpammer(bot *tgbotapi.BotAPI, start int64, dict *map[int]int, update *tgbotapi.Update, msg tgbotapi.MessageConfig) error {
	t := time.Now().UnixNano()
	elapsed := (t - start) / 1000000
	(*dict)[update.Message.From.ID]++
	if elapsed <= 5000 && (*dict)[update.Message.From.ID] > 5 {
		msg.Text = "You are spammer"
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
	} else if elapsed >= 5000 {
		start = time.Now().UnixNano()
		*dict = map[int]int{}
	}
	return nil
}
