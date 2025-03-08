package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

/*
*
fields are capitalized to mark them exported and accessible from outside the package
json tag to provide metadata, that will be useful to convert the entity into json when sending over network
*/
type Todo struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var DB *sql.DB

func initDatabase() {
	var err error
	// for sqlite, db is just a file on the disk
	DB, err = sql.Open("sqlite3", "./todo.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		status TEXT 
    )`)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateToDo(title string, status string) (int64, error) {
	result, err := DB.Exec("INSERT INTO todos (title, status) VALUES (?, ?)", title, status)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func DeleteToDo(id int64) error {
	_, err := DB.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}

func ReadToDoList() ([]Todo, error) {
	rows, err := DB.Query("SELECT id, title, status FROM todos")
	if err != nil {
		log.Println("error at reading todo list")
		return nil, err
	}
	defer rows.Close()

	// initialize list
	todos := make([]Todo, 0)

	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.Id, &todo.Title, &todo.Status)
		if err != nil {
			log.Println("error at reading todo list")
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}
