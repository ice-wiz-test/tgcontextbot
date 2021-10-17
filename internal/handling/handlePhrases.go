package handling

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func CheckMSG(firstPair *[]string, secondPair *[]string, listOfExceptions *[]string, newUpd tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	for i := 0; i < len(*firstPair); i++ {
		var h int64 = 0
		for j := 0; j < len(*listOfExceptions); j++ {
			if (*firstPair)[i] == (*listOfExceptions)[j] {
				h++
			}
		}
		if h == 0 {
			if strings.Contains(newUpd.Message.Text, (*firstPair)[i]) {
				msg := tgbotapi.NewMessage(newUpd.Message.Chat.ID, "")
				msg.Text = (*secondPair)[i]

				_, err := bot.Send(msg)
				if err != nil {
					HandleError(err)
				}
			}
		}

	}

	return nil
}
