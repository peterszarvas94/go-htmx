package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

	_ "github.com/libsql/libsql-client-go/libsql"
)

func db() (*sql.DB, error) {
	url, urlFound := os.LookupEnv("DB_URL")
	if !urlFound {
		return nil, errors.New("DB_URL not found")
	}

	token, token_found := os.LookupEnv("DB_TOKEN")
	if !token_found {
		return nil, errors.New("DB_TOKEN not found")
	}

	connectionStr := fmt.Sprintf("%s?authToken=%s", url, token)

	db, dbErr := sql.Open("libsql", connectionStr)
	if dbErr != nil {
		return nil, dbErr
	}

	return db, nil
}

func GetTodos() ([]TodoData, error) {
	db, dbErr := db()
	if dbErr != nil {
		return nil, dbErr
	}

	rows, queryErr := db.Query("SELECT * FROM todos")
	if queryErr != nil {
		return nil, queryErr
	}
	defer rows.Close()

	todos := []TodoData{}
	for rows.Next() {
		var id int
		var text string
		scanErr := rows.Scan(&id, &text)
		if scanErr != nil {
			return nil, scanErr
		}

		todo := TodoData{
			Id:   id,
			Text: text,
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func AddTodo(newText string) (TodoData, error) {
	db, dbErr := db()
	if dbErr != nil {
		return TodoData{}, dbErr
	}

	_, mutationErr := db.Exec("INSERT INTO todos (text) VALUES (?)", newText)
	if mutationErr != nil {
		return TodoData{}, mutationErr
	}

	query, queryErr := db.Query("SELECT * FROM todos WHERE text = ?", newText)
	if queryErr != nil {
		return TodoData{}, queryErr
	}

	var id int
	var text string

	for query.Next() {
		scanErr := query.Scan(&id, &text)
		if scanErr != nil {
			return TodoData{}, scanErr
		}
	}

	todo := TodoData{
		Id:   id,
		Text: text,
	}

	return todo, nil
}

func GetTodoById(id int) (TodoData, error) {
	db, dbErr := db()
	if dbErr != nil {
		return TodoData{}, dbErr
	}

	query, queryErr := db.Query("SELECT * FROM todos WHERE id = ?", id)
	if queryErr != nil {
		return TodoData{}, queryErr
	}

	var text string

	for query.Next() {
		scanErr := query.Scan(&id, &text)
		if scanErr != nil {
			return TodoData{}, scanErr
		}
	}

	todo := TodoData{
		Id:   id,
		Text: text,
	}

	return todo, nil
}

func DeleteTodoById(id int) error {
	db, dbErr := db()
	if dbErr != nil {
		return dbErr
	}

	_, mutationErr := db.Exec("DELETE FROM todos WHERE id = ?", id)
	if mutationErr != nil {
		return mutationErr
	}

	return nil
}

func AddUser(newUsername string, newEmail string, newPassword string) (UserData, error) {
	db, dbErr := db()
	if dbErr != nil {
		return UserData{}, dbErr
	}

	hashedPassword, hashErr := HashPassword(newPassword)
	if hashErr != nil {
		return UserData{}, hashErr
	}
	_, mutationErr := db.Exec(
		"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		newUsername, newEmail, hashedPassword,
	)
	if mutationErr != nil {
		return UserData{}, mutationErr
	}

	query, queryErr := db.Query("SELECT * FROM users WHERE username = ?", newUsername)
	if queryErr != nil {
		return UserData{}, queryErr
	}

	var id int
	var username string
	var email string
	var password string

	for query.Next() {
		scanErr := query.Scan(&id, &username, &email, &password)
		if scanErr != nil {
			return UserData{}, scanErr
		}
	}

	user := UserData{
		Id:       id,
		Username: username,
		Email:    email,
	}

	return user, nil
}

func LoginUser(usernameOrEmail string, password string) (UserData, error) {
	db, dbErr := db()
	if dbErr != nil {
		return UserData{}, dbErr
	}

	query, queryErr := db.Query(
		"SELECT id, username, email, password as hash FROM users WHERE username = ? OR email = ?",
		usernameOrEmail, usernameOrEmail,
	)
	if queryErr != nil {
		return UserData{}, queryErr
	}

	var id int
	var username string
	var email string
	var hash string

	for query.Next() {
		scanErr := query.Scan(&id, &username, &email, &hash)
		if scanErr != nil {
			return UserData{}, scanErr
		}
	}

	if id == 0 {
		return UserData{}, errors.New("User not found")
	}

	matchErr := CheckPasswordHash(hash, password)
	if matchErr != nil {
		return UserData{}, matchErr
	}

	user := UserData{
		Id:       id,
		Username: username,
		Email:    email,
	}

	return user, nil

}

func GetUserById(id string) (UserData, error) {
	db, dbErr := db()
	if dbErr != nil {
		return UserData{}, dbErr
	}

	query, queryErr := db.Query("SELECT * FROM users WHERE id = ?", id)
	if queryErr != nil {
		return UserData{}, queryErr
	}

	var username string
	var email string
	var password string

	for query.Next() {
		scanErr := query.Scan(&id, &username, &email, &password)
		if scanErr != nil {
			return UserData{}, scanErr
		}
	}

	idInt, idErr := strconv.Atoi(id)
	if idErr != nil {
		return UserData{}, idErr
	}

	user := UserData{
		Id:       idInt,
		Username: username,
		Email:    email,
	}

	return user, nil
}
