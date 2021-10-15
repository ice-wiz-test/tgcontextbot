package startup

import (
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

//TODO - we should standarize handling errors inside the BotCommandHandle function
func BotCommandHandle(newUpd tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	msg := tgbotapi.NewMessage(newUpd.Message.Chat.ID, "")
	switch newUpd.Message.Command() {
	case "start":
		msg.Text = "Добро пожаловать! Ознакомьтесь с доступными командами для данного бота: " +
			"\n /addchat - позволяет добавлять бота в новый чат \n /addblacklist - позволяет добавлять слова, употребление которых в чате нежелательно " +
			"\n /watchblacklist - позволяет просматривать добавленные в черный список слдова " +
			"\n /help - позволяет получить помощь"
	case "addchat":

		var ErrorWithHandlingNewChat error = nil
		ErrorWithHandlingNewChat = handle.BotNewChatHandle(newUpd, bot)

		if ErrorWithHandlingNewChat != nil {
			return ErrorWithHandlingNewChat
		}

		return nil

	case "getpairs":

		firstptr, secondptr, ErrWithHandling, answer := connect.GetAllPairsFromChat(newUpd.Message.Chat.ID)

		msg.Text = answer

		if ErrWithHandling != nil {
			log.Println(ErrWithHandling)
			return ErrWithHandling
		}
		if firstptr != nil && len(*firstptr) != 0 {
			for i := 0; i < len(*firstptr); i++ {
				msg.Text += "\n"
				msg.Text += (*firstptr)[i]
				msg.Text += " -> "
				msg.Text += (*secondptr)[i]
			}
		} else {
			msg.Text = "Пар нету"
		}
	case "deletesubstitute":

		ErrWithParse, s := connect.DeleteWordFromChat(newUpd.Message.Chat.ID, newUpd.Message.Text)

		if ErrWithParse != nil {
			log.Println(ErrWithParse)
		}

		msg.Text = s

	case "addblacklist":

		var id = newUpd.Message.Chat.ID

		var s string
		s = strings.Trim(newUpd.Message.Text, "/addblacklist")

		if len(s) >= 3 {
			err := connect.AddWordToBlacklist(id, s)
			if err != nil {
				msg.Text = "Что-то пошло не так. Проверьте, что ваш чат добавлен в нашу базу данных."
			} else {
				msg.Text = "Мы успешно добавили слова."
			}
		} else {
			msg.Text = "Чтобы бот не отвечал на практически все сообщения, нельзя банить слова меньше 3 букв. Ну изивните."
		}

	case "watchblacklist":
		allWords, errr := connect.GetAllBadWordsByChat(newUpd.Message.Chat.ID)

		if errr != nil {
			return errr
		}

		for i := 0; i < len(*allWords); i++ {
			msg.Text += (*allWords)[i]
			msg.Text += "\n"
		}

		if msg.Text == "" {
			msg.Text = "В этом чате еще нету слов, которые мы отслеживаем."
		}
	case "deletefromblacklist":
		var id = newUpd.Message.Chat.ID

		var s string
		s = strings.Trim(newUpd.Message.Text, "/deletefromblacklist")

		err := connect.DeleteWordFromBlacklist(id, s)

		if err != nil {
			log.Println(err)
			msg.Text = "Что-то пошло не так. Проверьте, что ваш чат добавлен в нашу базу данных."
		} else {
			msg.Text = "Либо мы успешно добавили слова, либо ваш чат не в базе данных. 50/50"
		}
	case "guide":
		msg.Text = "https://github.com/ice-wiz-test/tgcontextbot/blob/main/guide.md"
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}

	case "setsubstitutewith":

		Err, answer := connect.AddWordToID(newUpd.Message.Text, newUpd.Message.Chat.ID)

		if Err != nil {
			return Err
		}

		msg.Text = answer

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
			//TODO - this should definitely be a separate function
			t := time.Now().UnixNano()
			elapsed := (t - start) / 1000000
			dict[update.Message.From.ID]++
			if elapsed <= 5000 && dict[update.Message.From.ID] > 5 {
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
					log.Println(err)
				}
				// at the moment, this simply logs all the errors we encounter during running. I do not yet see a better way of hadnling this
			} else {

				allWords, errr := connect.GetAllBadWordsByChat(update.Message.Chat.ID)

				if errr != nil {
					log.Println(errr)
				} else {
					if allWords == nil {
					} else {
						if handle.CheckProf(allWords, update.Message.Text) {
							msg.Text = "Вы сказали запрещенное слово, не надо так."
							_, _ = bot.Send(msg)
						}

						firstptr, secondptr, err, _ := connect.GetAllPairsFromChat(update.Message.Chat.ID)

						if err == nil {
							err = handle.CheckMSG(firstptr, secondptr, update, bot)
							if err != nil {
								log.Println(err)
							}
						} else {
							log.Println(err)
						}
					}
				}

			}
		}
	}

	return nil
}
