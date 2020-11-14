package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var todos []Todo

func getTodoHandler(c *gin.Context) {
	status := c.Query("status")
	items := []Todo{}
	queryAll()
	for _, item := range todos {
		if status != "" {
			if item.Status == status {
				items = append(items, item)
				//e.g. http://localhost:1234/todos?status=completed
			}
		} else {
			items = append(items, item)
		}

	}
	c.JSON(http.StatusOK, items)
}

func queryAll() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
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
		//fmt.Println(id, title, status)
	}
	fmt.Println("query all todo success")
	//fmt.Printf("%v", todos)
}

func main() {
	r := gin.Default()

	r.GET("/todos", getTodoHandler)

	r.Run(":1234")
}
