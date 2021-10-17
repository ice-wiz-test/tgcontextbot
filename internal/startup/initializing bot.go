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

//TODO - we should standarize handling errors inside the BotCommandHandle function
func BotCommandHandle(newUpd tgbotapi.Update, bot *tgbotapi.BotAPI) error {

	msg := tgbotapi.NewMessage(newUpd.Message.Chat.ID, "")
	switch newUpd.Message.Command() {
	case "addbadwordexception":
		var s string = strings.TrimLeft(newUpd.Message.Text, "/addbadwordexception")
		s = strings.TrimSpace(s)
		var mass []string = strings.Split(s, "||")
		if len(mass) != 3 {
			msg.Text = "Неверный формат. Используйте формат excepted_phrase||excepted_username|| для таких запросов"
		} else {
			ErrWithDB, str := connect.AddExceptionToChat(newUpd.Message.Chat.ID, mass[1], mass[0])
			if ErrWithDB != nil {
				log.Println(ErrWithDB)

			}

			msg.Text = str
		}

	case "start":
		msg.Text = "Добро пожаловать! Ознакомьтесь с доступными командами для данного бота: \n" +
			"/start - команда, запускающая бота, выдает список доступных команд \n" +
			"/help - команда, выдающая список доступных команд\n" +
			"/guide - команда, возвращающая ссылку на гайд\n" +
			"/addchat + ID чата - команда, позволяющая добавить чат в базу данных. Нужна, чтобы бот запомнил фильтр запрещенных слов в данном чате.\n" +
			"/addblacklist + слова - команда, добавляющая запрещенные слова в базу данных. После этого бот будет сообщать, если данное слово было употреблено\n" +
			"/watchblacklist - команда, позволяющая просмотреть список запрещенных слов для данного чата\n" +
			"/deletefromblacklist + слова - команда, позволяющая удалить выбранные запрещенные слова для данного чата\n" +
			"/setsubstitutewit + слово + || + слово - команда, устанавливающая соответствие между двумя словами. В последствии если первое слово будет употреблено в чате бот вернет слово, с которым это соответствие было установлено. Для того чтобы бот работал корректно нужно ввести два слова и разделить их знаком '||'\n" +
			"/getpairs -  возвращает все слова когда-либо употребленные в чате, с которыми было установленно соответствие предыдущей командой"

	case "addexceptiontosubstitute":
		var s string = strings.TrimLeft(newUpd.Message.Text, "/addexceptiontosubstitute")
		s = strings.TrimSpace(s)
		var mass []string = strings.Split(s, "||")
		if len(mass) != 3 {
			msg.Text = "Неверный формат. Используйте формат excepted_phrase||excepted_username|| для таких запросов"
		} else {
			ErrWithDB, str := connect.AddException(newUpd.Message.Chat.ID, mass[0], mass[1])
			if ErrWithDB != nil {
				log.Println(ErrWithDB)

			}

			msg.Text = str
		}
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
	case "getexcepted":
		fmt.Println(newUpd.Message.Text)

		firstptr, secondptr, ErrWithDB, _ := connect.GetExceptions(newUpd.Message.Chat.ID)

		if ErrWithDB == nil {
			if len(*firstptr) == 0 {
				msg.Text = "В данном чате нет исключений."
			} else {
				for i := 0; i < len(*firstptr); i++ {
					msg.Text += (*firstptr)[i]
					msg.Text += " не действует на ->"
					msg.Text += (*secondptr)[i]
					msg.Text += "\n"
				}
			}
		} else {
			return ErrWithDB
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
		s = strings.TrimLeft(newUpd.Message.Text, "/addblacklist")

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

	case "setsubstitutewith":

		Err, answer := connect.AddWordToID(newUpd.Message.Text, newUpd.Message.Chat.ID)

		if Err != nil {
			return Err
		}

		msg.Text = answer

	case "help":
		msg.Text = "Добро пожаловать! Ознакомьтесь с доступными командами для данного бота: \n" +
			"/start - команда, запускающая бота, выдает список доступных команд \n" +
			"/help - команда, выдающая список доступных команд\n" +
			"/guide - команда, возвращающая ссылку на гайд\n" +
			"/addchat + ID чата - команда, позволяющая добавить чат в базу данных. Нужна, чтобы бот запомнил фильтр запрещенных слов в данном чате.\n" +
			"/addblacklist + слова - команда, добавляющая запрещенные слова в базу данных. После этого бот будет сообщать, если данное слово было употреблено\n" +
			"/watchblacklist - команда, позволяющая просмотреть список запрещенных слов для данного чата\n" +
			"/deletefromblacklist + слова - команда, позволяющая удалить выбранные запрещенные слова для данного чата\n" +
			"/setsubstitutewit + слово + || + слово - команда, устанавливающая соответствие между двумя словами. В последствии если первое слово будет употреблено в чате бот вернет слово, с которым это соответствие было установлено. Для того чтобы бот работал корректно нужно ввести два слова и разделить их знаком '||'\n" +
			"/getpairs -  возвращает все слова когда-либо употребленные в чате, с которыми было установленно соответствие предыдущей командой"

	case "deleteexception":
		fmt.Println(newUpd.Message.Text)
		var s string = strings.TrimLeft(newUpd.Message.Text, "/deletexcepetion")
		s = strings.TrimSpace(s)
		var mass []string = strings.Split(s, "||")

		if len(mass) != 3 {
			msg.Text = "Неверный формат. Используйте формат excepted_phrase||excepted_username|| для таких запросов"
		} else {
			ErrWithDB, str := connect.DeleteExceptedWord(newUpd.Message.Chat.ID, mass[1], mass[0])
			if ErrWithDB != nil {
				log.Println(ErrWithDB)

			}

			msg.Text = str
		}

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
			err := handle.FindSpammer(bot, start, &dict, update.Message.From.ID, msg)
			if err != nil {
				log.Println(err)
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
						strptr, err, _ := connect.GetExceptionsByUsername(update.Message.Chat.ID, update.Message.From.UserName)
						ptr, Err := connect.GetAllBadWordsByChat(update.Message.Chat.ID)
						if err == nil && Err == nil {
							if handle.CheckProf(ptr, update.Message.Text, strptr) {
								msg.Text = "Вы сказали запрещенное в данном чате слово."
								_, err = bot.Send(msg)

								if err != nil {
									log.Println(err)
								}
							}
						} else {
							if err != nil {
								log.Println(err)
							}

							if Err != nil {
								log.Println(Err)
							}
						}

						firstptr, secondptr, err, _ := connect.GetAllPairsFromChat(update.Message.Chat.ID)

						exceptPTR, ErrWithParse, _ := connect.GetWordsByException(update.Message.Chat.ID, update.Message.From.UserName)

						if err == nil && ErrWithParse == nil {
							err = handle.CheckMSG(firstptr, secondptr, exceptPTR, update, bot)
						} else {
							if err != nil {
								log.Println(err)
							}

							if ErrWithParse != nil {
								log.Println(ErrWithParse)
							}
						}
					}
				}
			}

		}
	}

	return nil
}
