package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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

func createTodoHandler(c *gin.Context) {
	t := Todo{}
	//r.body and read body -> bind json and send to &t

	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := insertTodo(t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, todo)

}

func getTodoByIdHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id")) //atomic to integer
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	todos, err := queryById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
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
		log.Fatal("Connect to database error", err)
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

func insertTodo(t Todo) ([]Todo, error) {

	var todos []Todo
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	row := db.QueryRow("INSERT INTO todos (title, status) values ($1, $2)  RETURNING id", t.Title, t.Status)
	var id int
	err = row.Scan(&id)
	if err != nil {
		return []Todo{}, fmt.Errorf("can't scan id %s", err)
	}
	todo := Todo{id, t.Title, t.Status}
	todos = append(todos, todo)
	fmt.Println("insert todo success id: ", id)
	return todos, nil
}

func queryById(rowId int) ([]Todo, error) {

	var todos []Todo
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, title, status FROM todos where id=$1")
	if err != nil {
		return []Todo{}, fmt.Errorf("can't prepare query one row statement", err)
	}

	row := stmt.QueryRow(rowId)
	var id int
	var title, status string

	err = row.Scan(&id, &title, &status)
	if err != nil {
		return []Todo{}, fmt.Errorf("can't scan row into variables %s", err)
	}

	fmt.Println("one row", id, title, status)
	todo := Todo{id, title, status}
	todos = append(todos, todo)
	return todos, nil
}

func main() {
	r := gin.Default()

	r.GET("/todos", getTodoHandler)
	r.GET("/todos/:id", getTodoByIdHandler)
	r.POST("/todos", createTodoHandler)
	r.Run(":1234")
}
