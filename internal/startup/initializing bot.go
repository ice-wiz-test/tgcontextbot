package startup

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	handle "tgcontextbot/internal/handling"
)

func InitializeBot() (error, *tgbotapi.BotAPI) {
	data, err := os.ReadFile("internal/startup/bottoken.txt");
	if(err != nil) {
		return err, nil
	}

	token := string(data)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err,nil
	}
	return nil, bot


}



func BotCommandHandle(newUpd tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	msg := tgbotapi.NewMessage(newUpd.Message.Chat.ID, "")
	switch newUpd.Message.Command() {
	case "start":
		msg.Text = "Добро пожаловать! Ознакомьтесь с доступными командами для данного бота"
	case "addchat":
		fmt.Println(newUpd.Message.Text)
		var ErrorWithHandlingNewChat error = nil
		ErrorWithHandlingNewChat = handle.BotNewChatHandle(newUpd, bot)

		if ErrorWithHandlingNewChat != nil {
			return ErrorWithHandlingNewChat
		}

		return nil

	default:
		msg.Text = "Я не знаю такой команды, простите"
	}
	if msg.Text == "" {
		msg.Text = "Временный костыль для тестирования"
	}
	_, err := bot.Send(msg)

	if err != nil {
		return err
	}

	return nil
}

func ServeBot(bot *tgbotapi.BotAPI) {
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {


		if update.Message.IsCommand() {
			err := BotCommandHandle(update, bot)

			if err != nil {
				log.Panic(err)
			}
			// in the future this should probably return the error directly to the main program so that we can actually handle it
		}
	}
}