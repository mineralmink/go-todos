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

func getTodoHandler(c *gin.Context) {
	todos, err := queryAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusOK, todos)
		return
	}

	items := []Todo{}
	for _, item := range todos {
		if item.Status == status {
			items = append(items, item)
			//e.g. http://localhost:1234/todos?status=completed
		}
	}
	c.JSON(http.StatusOK, items)
}

func queryAll() ([]Todo, error) {

	var todos []Todo
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	defer db.Close()

	stmt, err := db.Prepare("SELECT id,title, status FROM todos")
	if err != nil {
		return []Todo{}, fmt.Errorf("can't prepare query all todos statement %s", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return []Todo{}, fmt.Errorf("can't query all todos %s", err)
	}

	for rows.Next() {
		var id int
		var title, status string

		err := rows.Scan(&id, &title, &status)
		if err != nil {
			return []Todo{}, fmt.Errorf("can't Scan row into variable %s", err)
		}
		todo := Todo{id, title, status}
		todos = append(todos, todo)
	}
	fmt.Println("query all todo success")
	return todos, nil
}

func main() {
	r := gin.Default()

	r.GET("/todos", getTodoHandler)

	r.Run(":1234")
}
