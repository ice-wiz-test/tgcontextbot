package storage

import (
	"context"
	db "github.com/jackc/pgx/v4"
	"strings"
)

func AddWordToBlacklist(idd int64, badWord string) error {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")

	if err != nil {
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
		flag = false
		for j = 0; j < len(badWordByChat); j++ {
			if allBadWords[i] == badWordByChat[j] {
				flag = true
			}
		}

		if !flag {
			_, err = conn.Exec(context.Background(), "insert into added_to_chats (id, badword) values ($1, $2)", idd, allBadWords[i])

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GetAllBadWordsByChat(idd int64) (*[]string, error) {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")

	if err != nil {
		return nil, err
	}

	defer conn.Close(context.Background())

	newRows, Err := conn.Query(context.Background(), "select badword from added_to_chats where id = $1", idd)

	if Err != nil {
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

	return &allWords, nil

}

func DeleteWordFromBlacklist(idd int64, badWord string) error {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")

	if err != nil {
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

	var flag = false
	var i = 0
	var j = 0
	for i = 0; i < len(allBadWords); i++ {
		flag = true
		for j = 0; j < len(badWordByChat); j++ {
			if allBadWords[i] == badWordByChat[j] {
				flag = false
			}
		}

		if !flag {
			_, err = conn.Exec(context.Background(), "delete from added_to_chats where id = $1 and badword = $2", idd, allBadWords[i])

			if err != nil {
				return err
			}
		}
	}

	return nil
}
