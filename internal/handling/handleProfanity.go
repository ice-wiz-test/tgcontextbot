package handling

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func BotHandleProfanity(newUpd tgbotapi.Update, bot *tgbotapi.BotAPI) ([]string, error) {
	msg := tgbotapi.NewMessage(newUpd.Message.Chat.ID, "")

	TextOfMessage := newUpd.Message.Text
	TextOfMessage = strings.ReplaceAll(TextOfMessage, "/addblacklist", "")
	fmt.Println(TextOfMessage)
	var profanity []string = nil
	profanity = strings.Fields(TextOfMessage)
	//fmt.Println(profanity)
	profanity = UniqueNonEmptyElementsOf(profanity)

	fmt.Println(profanity)
	msg.Text = "Ok"

	_, err := bot.Send(msg)

	if err != nil {
		return profanity, err
	}

	return profanity, err
}

func FindProfanity(profanity []string, newUpd tgbotapi.Update, bot *tgbotapi.BotAPI) (bool, error) {
	msg := tgbotapi.NewMessage(newUpd.Message.Chat.ID, "")

	TextOfMessage := newUpd.Message.Text

	words := strings.Split(TextOfMessage, " ")

	for i := 0;i < len(words);i++ {
		for j := 0; j < len(profanity); j++ {
			if words[i] == profanity[j] {
				var err error = nil
				msg.Text = "Curse word have been found"
				_, err = bot.Send(msg)
				if err != nil {
					return false, err
				}
				return true, err
			}
		}
	}
	var err error = nil
	return false,err
}

func UniqueNonEmptyElementsOf(s []string) []string {
	unique := make(map[string]bool, len(s))
	us := make([]string, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				us = append(us, elem)
				unique[elem] = true
			}
		}
	}

	return us
}
