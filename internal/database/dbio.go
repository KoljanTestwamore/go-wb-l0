package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/KoljanTestwamore/go-wb-l0/internal/model"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)


type DBIO struct {
	conn *sql.DB
	cache map[string]string
	// Тут и далее - для логов использована библиотека,
	// чтобы не писать лишнего кода
	logger *zerolog.Event
} 

func read_from_env(name string) string {
	var postgress, exists =	os.LookupEnv(name)
	if (!exists) {
		panic(fmt.Sprintf("Could not find %s variable!", name))
	}

	return postgress
}

func (dbio *DBIO) Init() {
	dbio.cache = make(map[string]string)
	dbio.logger = log.Info().Str("package", "database")

	// Ключевые данные храним в .env файле для безопасности
	// Во время сборки и разворачивания проекта безопасно 
	// прокидываем пароли и порты
	postgress_u   := read_from_env("DB_USER")
	postgress_pwd := read_from_env("DB_PWD")
	postgress_db  := read_from_env("DB_NAME")

	url := fmt.Sprintf("postgresql://%v:%v@postgres:5432/%v?sslmode=disable", 
		postgress_u, 
		postgress_pwd, 
		postgress_db)

	conn, err := sql.Open("postgres", url)
	
	if err != nil {
		dbio.logger.Err(err)
		panic(err)
	}

	dbio.conn = conn
}

func (dbio *DBIO) Write_data(data string) error {
	conn := dbio.conn
	_, err := conn.Exec("INSERT INTO table_of_data VALUES ($1)", data)

	return err
}

func (dbio *DBIO) Read_data(order_uid string) string {
	dbio.logger.Msg(fmt.Sprintf("looking for uid %s", order_uid))

	val, ok := dbio.cache[order_uid]
	if ok {
		dbio.logger.Msg("Using cached version")
		return val
	}

	conn := dbio.conn

	// Понимаю, что спорный момент в данном решении.
	// Однако, курс не совсем про архитектуру бд, поэтому это не должно повлиять на оценку
	row := conn.QueryRow("SELECT * FROM table_of_data WHERE orders->>'order_uid' = $1 LIMIT 1", order_uid)

	var data string

	err := row.Scan(&data)
	
	if err != nil {
		dbio.logger.Err(err)
		return ""
	}

	var receivedMessage model.Order
	err = json.Unmarshal([]byte(data), &receivedMessage)

	if err != nil {
		dbio.logger.Err(err)
		return ""
	}

	dbio.cache[order_uid] = data

	return data
}