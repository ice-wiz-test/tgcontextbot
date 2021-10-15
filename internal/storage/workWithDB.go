package storage

import (
	"context"
	db "github.com/jackc/pgx/v4"
	"log"
)

//TODO - rewrite all functions that work with databases into (error, string) format
//TODO - rewrite this TODO into the correct file, i am too tired to do that now
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
		return err
	}

	defer conn.Close(context.Background())

	_, Err := conn.Exec(context.Background(), "insert into added_to_chats (id, badword) values ($1, 'placeholder')", idd)
	if Err != nil {
		return Err
	}

	return nil
}
