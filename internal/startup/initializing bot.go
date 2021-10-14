package startup

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strings"
	handle "tgcontextbot/internal/handling"
	connect "tgcontextbot/internal/storage"
	"time"
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

	start := time.Now().UnixNano()
	dict := map[int]int{}
	for update := range updates {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if update.Message != nil {
			t := time.Now().UnixNano()
			elapsed := (t - start) / 1000000
			dict[update.Message.From.ID]++
			fmt.Println(dict[update.Message.From.ID])
			if elapsed >= 5000 && dict[update.Message.From.ID] > 30 {
				msg.Text = "You are spammer"
				_, err := bot.Send(msg)
				if err != nil {
					return err
				}
			} else if elapsed >= 5000 {
				start = time.Now().UnixNano()
				dict = map[int]int{}
			}
			if update.Message.IsCommand() {

				err := BotCommandHandle(update, bot)

				if err != nil {
					log.Panic(err)
				}
				// in the future this should probably return the error directly to the main program so that we can actually handle it
			} else {
				fmt.Println("KOSTYL")

				allWords, errr := connect.GetAllBadWordsByChat(update.Message.Chat.ID)

				if errr != nil {
					log.Println(errr)
				} else {
					if allWords == nil {
						log.Println("This chat does not have any words")
					} else {
						if handle.CheckProf(&allWords, update.Message.Text) {
							msg.Text = "Вы сказали запрещенно слово, не надо так."
							_, _ = bot.Send(msg)
						}
					}
				}

			}
		}
	}

	return nil
}
