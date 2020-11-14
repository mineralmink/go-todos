package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Todo struct {
	ID     int
	Title  string
	Status string
}

var todos []Todo

func main() {
	db, err := sql.Open("postgres", "postgres://jeoysqgj:Yqpkng49GujaIUn9LrGfzP1bRD3JHFAM@suleiman.db.elephantsql.com:5432/jeoysqgj")
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	defer db.Close()

	stmt, err := db.Prepare("SELECT id,title, status FROM todos")
	if err != nil {
		log.Fatal("can't prepare query all todos statement", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("can't query all todos", err)
	}

	for rows.Next() {
		var id int
		var title, status string

		err := rows.Scan(&id, &title, &status)
		if err != nil {
			log.Fatal("can't Scan row into variable", err)
		}
		todo := Todo{id, title, status}
		todos = append(todos, todo)
		fmt.Println(id, title, status)
	}
	fmt.Println("query all todo success")
	fmt.Printf("%v", todos)
}
