package startup

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	handle "tgcontextbot/internal/handling"
)

func InitializeBot() (error, *tgbotapi.BotAPI) {
	data, err := os.ReadFile("internal/startup/bottoken.txt")
	if err != nil {
		return err, nil
	}

	token := string(data)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err, nil
	}
	return nil, bot

}

func BotCommandHandle(newUpd tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	msg := tgbotapi.NewMessage(newUpd.Message.Chat.ID, "")
	switch newUpd.Message.Command() {
	case "start":
		msg.Text = "Добро пожаловать! Ознакомьтесь с доступными командами для данного бота: " +
			"\n /addchat - позволяет добавлять бота в новый чат \n /addblacklist - позволяет добавлять слова, употребление которых в чате нежелательно " +
			"\n /watchblacklist - позволяет просматривать добавленные в черный список слдова " +
			"\n /help - позволяет получить помощь"
	case "addchat":
		fmt.Println(newUpd.Message.Text)
		var ErrorWithHandlingNewChat error = nil
		ErrorWithHandlingNewChat = handle.BotNewChatHandle(newUpd, bot)

		if ErrorWithHandlingNewChat != nil {
			return ErrorWithHandlingNewChat
		}

		return nil
	case "addblacklist":
		fmt.Println(newUpd.Message.Text)
		var ErrorWithHandlingBlackList error = nil
		profanity,ErrorWithHandlingBlackList = handle.BotHandleProfanity(newUpd, bot)

		if ErrorWithHandlingBlackList != nil {
			return ErrorWithHandlingBlackList
		}

		return nil
	case "watchblacklist":
		fmt.Println(newUpd.Message.Text)
		fmt.Println(profanity)
		
		for i := 0;i < len(profanity);i++{
			msg.Text += profanity[i]
			msg.Text += "\n"
		}
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}

		return nil
	case "help":
		msg.Text = "Добро пожаловать! Ознакомьтесь с доступными командами для данного бота: " +
			"\n /addchat - позволяет добавлять бота в новый чат " +
			"\n /addblacklist - позволяет добавлять слова, употребление которых в чате нежелательно " +
			"\n /watchblacklist - позволяет просматривать добавленные в черный список слдова" +
			"\n /help - позволяет получить помощь"
	default:
		msg.Text = "Я не знаю такой команды, простите"
	}
	isProfanity, errr := handle.FindProfanity(profanity, newUpd, bot)
	if errr != nil{
		return errr
	}
	if isProfanity {
		msg.Text = "@" + newUpd.Message.From.UserName + " You have said curse word"
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

		fmt.Println("BREAKPOINT 1")

		if update.Message != nil {
			if update.Message.IsCommand() {
				fmt.Println("HELP")

				err := BotCommandHandle(update, bot)

				if err != nil {
					log.Panic(err)
				}
				// in the future this should probably return the error directly to the main program so that we can actually handle it
			} else {
				fmt.Println("Help")
			}
		}
	}


		// TODO - this should handle non-command messages


}
