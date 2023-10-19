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

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS user_request (id INTEGER PRIMARY KEY, url TEXT, method TEXT, body TEXT, body_format TEXT, created_at DATETIME, updated_at DATETIME)")

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

func ReadDataById(id int, db *sql.DB) (*sql.Rows, error) {

	data, err := db.Query("SELECT * FROM user_request WHERE id = ?", id)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return data, nil
}

func ReadData(db *sql.DB) (*sql.Rows, error) {

	rows, err := db.Query("SELECT * FROM user_request")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return rows, nil
}

func UpdateData(id int, data *RequestData, db *sql.DB) (sql.Result, error) {
	// TODO -> Just updating the first row, I need to update the row with the id passed as parameter.
	// but for now I will just update the first row.

	result, err := db.Exec("UPDATE user_request SET url = ?, method = ?, body = ?, updated_at = datetime('now') WHERE id = 1", &data.Url, &data.Method, &data.Body)

	if err != nil {
		panic(err)
	}

	fmt.Println(result)
	return result, nil
}
