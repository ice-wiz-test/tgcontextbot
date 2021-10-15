package handling

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

func CheckProf(badWords *[]string, find string) bool {
	for i := 0; i < len(*(badWords)); i++ {
		if strings.Contains(find, (*badWords)[i]) {
			return true
		}
	}
	return false
}

func CheckMSG(firstPair *[]string, secondPair *[]string, newUpd tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	for i := 0; i < len(*firstPair); i++ {
		if strings.Contains(newUpd.Message.Text, (*firstPair)[i]) {
			msg := tgbotapi.NewMessage(newUpd.Message.Chat.ID, "")
			msg.Text = (*secondPair)[i]

			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
				return err
			}
		}

	}

	return nil
}
