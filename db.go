package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/libsql/libsql-client-go/libsql"
)

func db() *sql.DB {
	url, url_err := os.LookupEnv("DB_URL")
	if url_err != true {
		fmt.Println("Error:", url_err)
		os.Exit(1)
	}

	token, token_err := os.LookupEnv("DB_TOKEN")
	if token_err != true {
		fmt.Println("Error:", token_err)
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

func GetTodoById (id int) Todo {
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

func DeleteTodoById (id int) {
	db := db()

	_, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
