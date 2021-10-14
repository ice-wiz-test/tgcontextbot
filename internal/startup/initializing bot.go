package startup

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strings"
	handle "tgcontextbot/internal/handling"
	connect "tgcontextbot/internal/storage"
)

var profanity []string = nil

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
		/*
			var ErrorWithHandlingBlackList error = nil

			profanity, ErrorWithHandlingBlackList = handle.BotHandleProfanity(newUpd, bot)

			if ErrorWithHandlingBlackList != nil {
				return ErrorWithHandlingBlackList
			}

			return nil
		*/

		var id int64 = newUpd.Message.Chat.ID

		var s string
		s = strings.Trim(newUpd.Message.Text, "/addblacklist")
		fmt.Println(s)

		err := connect.AddWordToBlacklist(id, s)

		if err != nil {
			log.Println(err)
			msg.Text = "Что-то пошло не так. Проверьте, что ваш чат добавлен в нашу базу данных."
		} else {
			msg.Text = "Либо мы успешно добавили слова, либо ваш чат не в базе данных. 50/50"
		}

	case "watchblacklist":
		fmt.Println(newUpd.Message.Text)
		allWords, errr := connect.GetAllBadWordsByChat(newUpd.Message.Chat.ID)

		if errr != nil {
			return errr
		}

		for i := 0; i < len(allWords); i++ {
			msg.Text += allWords[i]
			msg.Text += "\n"
		}

		if msg.Text == "" {
			msg.Text = "В этом чате еще нету слов, которые мы отслеживаем."
		}

	case "help":
		msg.Text = "Добро пожаловать! Ознакомьтесь с доступными командами для данного бота: " +
			"\n /addchat - позволяет добавлять бота в новый чат " +
			"\n /addblacklist - позволяет добавлять слова, употребление которых в чате нежелательно " +
			"\n /watchblacklist - позволяет просматривать добавленные в черный список слдова" +
			"\n /help - позволяет получить помощь"
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

func ServeBot(bot *tgbotapi.BotAPI) error {
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {

		if update.Message != nil {
			if update.Message.IsCommand() {

				err := BotCommandHandle(update, bot)

				if err != nil {
					log.Panic(err)
				}
				// in the future this should probably return the error directly to the main program so that we can actually handle it
			} else {
				//fmt.Println("Help")
				/*msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				isProfanity, errr := handle.FindProfanity(profanity, update, bot)
				if errr != nil {
					return errr
				}
				//fmt.Println(isProfanity)
				if isProfanity == true {
					msg.Text = "@" + update.Message.From.UserName + " You have said curse word"
					_, err := bot.Send(msg)
					if err != nil {
						return err
					}
				}

				*/
			}
		}
	}

	return nil
}
