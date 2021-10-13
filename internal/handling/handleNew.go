package handling

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	stor "tgcontextbot/internal/storage"
)

func BotNewChatHandle(newUpd tgbotapi.Update, bot *tgbotapi.BotAPI)  error {

	msg := tgbotapi.NewMessage(newUpd.Message.Chat.ID, "")

	var textOfMessage string = ""
	textOfMessage = newUpd.Message.Text
	fmt.Println(textOfMessage)
	textOfMessage = strings.Trim(textOfMessage, "/addchat")
	textOfMessage = strings.TrimSpace(textOfMessage)
	var id int64 = 0
	var err error = nil
	id, err = strconv.ParseInt(textOfMessage, 10, 64)
	fmt.Println(textOfMessage)
	if err != nil {
		msg.Text = "Сообщение не является числом"
		_, Err := bot.Send(msg)
		if Err != nil {
			log.Println(Err)
			return Err
		}

		return err
		log.Println(err)
	}

	/*if chat.IsInChatDir(id) {
		msg.Text = "Чат уже добавлен."
	} else {
		chat.AddToChatDir(id)
		msg.Text = "Чат добавлен во внутреннюю базу данных"
	}
	*/

	if stor.CheckIfPresentInChats(id) {
		msg.Text = "Чат уже добавлен в базу данных!"
	} else {
		msg.Text = "Бот работает"
	}

	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}


