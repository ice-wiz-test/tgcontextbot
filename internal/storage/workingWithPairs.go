package storage

import (
	"context"
	db "github.com/jackc/pgx/v4"
	"strings"
	handle "tgcontextbot/internal/handling"
)

func AddWordToID(keyString string, idd int64) (error, string) {
	var s string
	s = strings.Trim(keyString, "/setsubstitutewith ")
	s = strings.TrimSpace(s)
	var allString []string

	allString = strings.Split(s, "||")

	if len(allString) != 3 {
		return nil, "В команде не два слова, или они разделены неправильными символами(не ||). Пример - /setsubstitutewith aboba||amongus||"
	}
	var first string
	var second string
	first = allString[0]
	second = allString[1]
	var checker string
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")

	if err != nil {
		handle.HandleError(err)
		return err, "Ошибка при соединении с базой данных"
	}
	defer func(conn *db.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			handle.HandleError(err)
		}
	}(conn, context.Background())

	newRows, Err := conn.Query(context.Background(), "select replace_phrase from chat_phrases where chat_id = $1 and find_phrase = $2", idd, first)

	if Err != nil {
		handle.HandleError(Err)
		return Err, "Мы не сумели подключиться к базе данных. Повторите запрос через какое-то время."
	}
	var cnt int64 = 0
	for newRows.Next() {
		cnt++
		ErrWithParse := newRows.Scan(&checker)
		if ErrWithParse != nil {
			return ErrWithParse, "Ошибка при работе с базой данных."
		}
	}
	if cnt != 0 {
		return nil, "Уже есть заменитель на эту фразу"
	}
	_, Err = conn.Exec(context.Background(), "insert into chat_phrases (chat_id, find_phrase, replace_phrase) values ($1, $2, $3)", idd, first, second)

	if Err != nil {
		handle.HandleError(Err)
		return Err, "Мы не сумели подключиться к базе данных. Повторите запрос через какое-то время."
	}

	return nil, "Набор фраз добавлен"
}

func GetAllPairsFromChat(idd int64) (*[]string, *[]string, error, string) {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")

	if err != nil {
		handle.HandleError(err)
		return nil, nil, err, "Мы не сумели установить соединение с базой данных"
	}
	defer func(conn *db.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			handle.HandleError(err)
		}
	}(conn, context.Background())
	newRows, err := conn.Query(context.Background(), "select find_phrase, replace_phrase from chat_phrases where chat_id = $1", idd)
	var firstPair []string = nil
	var secondPair []string = nil
	var firstWordInPair string
	var secondWordInPair string

	if newRows == nil {
		handle.HandleError(err)
		return nil, nil, nil, "В данном чате нету замен"
	}
	for newRows.Next() {
		ErrWithParse := newRows.Scan(&firstWordInPair, &secondWordInPair)

		if ErrWithParse != nil {
			handle.HandleError(err)
			return nil, nil, ErrWithParse, "Ошибка при работе с базой данных"
		}

		firstPair = append(firstPair, firstWordInPair)

		secondPair = append(secondPair, secondWordInPair)
	}

	return &firstPair, &secondPair, nil, "Все успешно."
}

func DeleteWordFromChat(idd int64, key string) (error, string) {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")
	var s = strings.TrimLeft(key, "/deletesubstitute")
	s = strings.TrimSpace(s)
	if err != nil {
		return err, "Мы не сумели установить соединение с базой данных"
	}
	defer func(conn *db.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			handle.HandleError(err)
		}
	}(conn, context.Background())

	_, err = conn.Exec(context.Background(), "delete from chat_phrases where chat_id = $1 and find_phrase = $2", idd, s)

	if err != nil {
		handle.HandleError(err)
		return err, "Что-то пошло не так. Попробуйте позже."
	} else {
		return nil, "Все успешно."
	}
}

func AddException(chatId int64, key string, excepted string) (error, string) {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")
	if err != nil {
		handle.HandleError(err)
		return err, "Мы не сумели установить соединение с базой данных"
	}
	defer func(conn *db.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			handle.HandleError(err)
		}
	}(conn, context.Background())

	rows, err := conn.Query(context.Background(), "select autoinc_id where phrase = $1 and id_of_excepted = $2 and chat_id= $3", key, excepted, chatId)
	var cnt int64
	for rows.Next() && cnt <= 2 {
		cnt++
	}
	if cnt != 0 {
		return nil, "Уже добавлен в исключения."
	}

	_, err = conn.Exec(context.Background(), "insert into exceptions (phrase, id_of_excepted, chat_id) values ($1, $2, $3)", key, excepted, chatId)

	if err != nil {
		handle.HandleError(err)
		return err, "Мы не сумели подключиться к базе данных."
	}

	return nil, "Мы добавили человека в список исключений"

}

func GetWordsByException(chatId int64, excepted string) (*[]string, error, string) {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")
	if err != nil {
		handle.HandleError(err)
		return nil, err, "Мы не сумели установить соединение с базой данных"
	}
	defer func(conn *db.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			handle.HandleError(err)
		}
	}(conn, context.Background())

	newRows, err := conn.Query(context.Background(), "select phrase from exceptions where chat_id = $1 and id_of_excepted = $2", chatId, excepted)

	var ret []string

	for newRows.Next() {
		var str string
		Err := newRows.Scan(&str)
		if Err != nil {
			handle.HandleError(Err)
			return nil, Err, "Ошибка в базе данных. Пожалуйста, повторите запрос через некоторое время."
		}
		ret = append(ret, str)
	}

	return &ret, nil, "Вот ваш список"
}

func DeleteExceptedWord(chatId int64, excepted string, key string) (error, string) {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")
	if err != nil {
		return err, "Мы не сумели установить соединение с базой данных"
	}
	defer func(conn *db.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			handle.HandleError(err)
		}
	}(conn, context.Background())

	_, ErrWithDel := conn.Exec(context.Background(), "delete from exceptions where chat_id = $1 and phrase = $2 and id_of_excepted = $3", chatId, key, excepted)

	if ErrWithDel != nil {
		handle.HandleError(ErrWithDel)
		return ErrWithDel, "Произошла ошибка при удалении, пожалуйста, повторите запрос в ближайшее время."
	}

	return nil, "Успешно удален из списка исключений."
}

func GetExceptions(chatId int64) (*[]string, *[]string, error, string) {
	conn, err := db.Connect(context.Background(), "postgres://postgres:password@localhost:5432/test")
	if err != nil {
		handle.HandleError(err)
		return nil, nil, err, "Мы не сумели установить соединение с базой данных"
	}
	defer func(conn *db.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			handle.HandleError(err)
		}
	}(conn, context.Background())

	newRows, Err := conn.Query(context.Background(), "select phrase, id_of_excepted from exceptions where chat_id = $1", chatId)
	if Err != nil {
		handle.HandleError(Err)
		return nil, nil, Err, "На сервере произошла ошибка. Повторите запрос позже."
	}
	var str1 string
	var str2 string

	var ret1 []string
	var ret2 []string

	for newRows.Next() {
		Err = newRows.Scan(&str1, &str2)

		if Err != nil {
			handle.HandleError(Err)
			return nil, nil, Err, "На сервере произошла ошибка. Повторите запрос позже."
		}

		ret1 = append(ret1, str1)
		ret2 = append(ret2, str2)
	}

	return &ret1, &ret2, nil, "Все успешно."

}
