package storage

import (
	"context"
	"fmt"
	db "github.com/jackc/pgx/v4"
	"log"
	"strings"
)

func CheckIfPresentInChats(idd int64) bool {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")
	if err != nil {
		log.Println(err)
		return false
	}

	defer conn.Close(context.Background())

	var chat_id int64

	err = conn.QueryRow(context.Background(), "select autoinc_id from added_to_chats where id = $1", idd).Scan(&chat_id)

	if err != nil {
		log.Println(err)
		return false
	}

	return true

}

func AddChatIDToDatabase(idd int64) error {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")
	if err != nil {
		log.Println(err)
		return err
	}

	defer conn.Close(context.Background())

	_, Err := conn.Exec(context.Background(), "insert into added_to_chats (id, badword) values ($1, 'placeholder')", idd)
	if Err != nil {
		log.Println(Err)
		return Err
	}

	return nil
}

func AddWordToBlacklist(idd int64, badWord string) error {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")

	if err != nil {
		log.Println(err)
		return err
	}

	defer conn.Close(context.Background())

	var allBadWords []string = strings.Split(badWord, " ")

	var badWordByChat []string = nil

	newRows, Err := conn.Query(context.Background(), "select badword from added_to_chats where id = $1", idd)

	if Err != nil {
		return Err
	}

	if newRows == nil {
		return nil
	}

	for newRows.Next() {
		var txt string
		errWithParse := newRows.Scan(&txt)
		if errWithParse != nil {
			return errWithParse
		}
		badWordByChat = append(badWordByChat, txt)
	}

	var flag bool = false
	var i int = 0
	var j int = 0
	for i = 0; i < len(allBadWords); i++ {
		fmt.Println(allBadWords[i])
		flag = false
		for j = 0; j < len(badWordByChat); j++ {
			if allBadWords[i] == badWordByChat[j] {
				flag = true
			}
		}

		if !flag {
			_, err = conn.Exec(context.Background(), "insert into added_to_chats (id, badword) values ($1, $2)", idd, allBadWords[i])

			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	return nil
}

func GetAllBadWordsByChat(idd int64) ([]string, error) {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer conn.Close(context.Background())

	newRows, Err := conn.Query(context.Background(), "select badword from added_to_chats where id = $1", idd)

	if Err != nil {
		log.Println(Err)
		return nil, Err
	}

	if newRows == nil {
		return nil, nil
	}

	var allWords []string = nil

	for newRows.Next() {
		var txt string

		errWithParse := newRows.Scan(&txt)

		if errWithParse != nil {
			return nil, errWithParse
		}

		allWords = append(allWords, txt)

	}

	return allWords, nil

}

/*func AddWordToID(keystring string, idd int64) (error, string) {
	var stringToSplit = strings.Trim(keystring, "/addphrase")
	stringToSplit = strings.TrimSpace(stringToSplit)

	var use []string = strings.Split(stringToSplit, "||")

	if len(use) < 2 {
		return nil, "В команде нету двух строк!"
	}


}

*/
