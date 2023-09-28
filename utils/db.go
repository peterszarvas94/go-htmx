package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/libsql/libsql-client-go/libsql"
)

func db() *sql.DB {
	url, url_found := os.LookupEnv("DB_URL")
	if !url_found {
		fmt.Println("Error:", url_found)
		os.Exit(1)
	}

	token, token_found := os.LookupEnv("DB_TOKEN")
	if !token_found {
		fmt.Println("Error:", token_found)
		os.Exit(1)
	}

	connectionStr := fmt.Sprintf("%s?authToken=%s", url, token)

	db, db_err := sql.Open("libsql", connectionStr)
	if db_err != nil {
		fmt.Println("Error:", db_err)
		os.Exit(1)
	}

	return db
}

func GetTodos() []Todo {
	db := db()

	rows, err := db.Query("SELECT * FROM todos")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var id int
		var text string
		err = rows.Scan(&id, &text)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		todo := Todo{
			Id:   id,
			Text: text,
		}
		todos = append(todos, todo)
	}

	return todos
}

func AddTodo(new_text string) Todo {
	db := db()

	_, err := db.Exec("INSERT INTO todos (text) VALUES (?)", new_text)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	query, err := db.Query("SELECT * FROM todos WHERE text = ?", new_text)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	var id int
	var text string

	for query.Next() {
		err = query.Scan(&id, &text)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}

	todo := Todo{
		Id:   id,
		Text: text,
	}

	return todo
}

func GetTodoById(id int) Todo {
	db := db()

	query, err := db.Query("SELECT * FROM todos WHERE id = ?", id)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	var text string

	for query.Next() {
		err = query.Scan(&id, &text)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}

	todo := Todo{
		Id:   id,
		Text: text,
	}

	return todo
}

func DeleteTodoById(id int) {
	db := db()

	_, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func NewUser(username string, email string, password string) (bool, error) {
	db := db()

	hashedPassword, hashErr := HashPassword(password)
	if hashErr != nil {
		return false, hashErr
	}

	_, dbErr := db.Exec(
		"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		username, email, hashedPassword,
	)
	if dbErr != nil {
		return false, dbErr
	}

	return true, nil
}

func GetUserByUsernameOrEmail(usernameOrEmail string, password string) (bool, error) {
	db := db()

	query, queryErr := db.Query(
		"SELECT id, username, email, password as hash FROM users WHERE username = ? OR email = ?",
		usernameOrEmail, usernameOrEmail,
	)
	if queryErr != nil {
		return false, queryErr
	}

	var id int
	var username string
	var email string
	var hash string

	for query.Next() {
		scanErr := query.Scan(&id, &username, &email, &hash)
		if scanErr!= nil {
			return false, scanErr
		}
	}

	notFoundErr := errors.New("User not found")
	if id == 0 {
		return false, notFoundErr
	}

	match, matchErr := CheckPasswordHash(hash, password)

	if match {
		return true, nil
	}
	return false, matchErr
}
