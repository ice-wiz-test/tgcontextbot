package main

import (
	"fmt"
	start "tgcontextbot/internal/startup"
)

func main() {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	/*
		conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: chat%v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())

		var chat_id int64
		err = conn.QueryRow(context.Background(), "select chat_id_index from chat where id = $1", 1).Scan(&chat_id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(chat_id)

	*/
	err, bot := start.InitializeBot()

	if err != nil {
		fmt.Println("Unable to start up, terminating.")
		return
	}

	fmt.Println("Started up")

	errr := start.ServeBot(bot)
	if errr != nil {
		return
	}

}
