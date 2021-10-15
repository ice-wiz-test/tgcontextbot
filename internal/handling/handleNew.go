package handling

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	stor "tgcontextbot/internal/storage"
)

func BotNewChatHandle(newUpd tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	msg := tgbotapi.NewMessage(newUpd.Message.Chat.ID, "")

	var textOfMessage string = ""
	textOfMessage = newUpd.Message.Text

	textOfMessage = strings.Trim(textOfMessage, "/addchat")
	textOfMessage = strings.TrimSpace(textOfMessage)
	var id int64 = 0
	id = newUpd.Message.Chat.ID

	var err error
	if stor.CheckIfPresentInChats(id) {
		msg.Text = "Чат уже добавлен в базу данных!"
	} else {
		check := stor.AddChatIDToDatabase(id)
		if check == nil {
			msg.Text = "Мы добавили ваш чат в базу данных."
		} else {
			return check
			msg.Text = "Что-то пошло не так("
		}
	}

	_, err = bot.Send(msg)
	if err != nil {
		return err
	}

	return nil

}
