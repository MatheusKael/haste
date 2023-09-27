package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type RequestData struct {
	Url         string
	Method      string
	Body        string
	Body_format string
}

// TODO -> I need to update the existing data, not just insert new data.

func ConnectAndCreateTable() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database.db")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res, err := db.Exec("CREATE TABLE IF NOT EXISTS  user_request (id INTEGER PRIMARY KEY, url TEXT, method TEXT, body TEXT, body_format TEXT, created_at DATETIME, updated_at DATETIME)")

	fmt.Println(res)

	return db, err
}

func InsertData(data *RequestData, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO user_request (url, method, body, body_format, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))", &data.Url, &data.Method, &data.Body, &data.Body_format)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func ReadData(db *sql.DB) (*sql.Rows, error) {

	rows, err := db.Query("SELECT * FROM user_request")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return rows, nil
}
