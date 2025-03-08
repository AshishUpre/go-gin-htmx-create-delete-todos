package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func main() {
	initDatabase()
	// defer closing of DB conn till end of main()
	defer DB.Close()

	e := gin.Default()

	// registering the templates directory
	e.LoadHTMLGlob("templates/*")

	// the gin context is context of curr http req, its passed as pointer as if its passed as struct
	// it maybe large -> large value to be put on stack
	// instead pointer => const size that points to the addr of the struct
	e.GET("/", func(c *gin.Context) {
		todos, err := ReadToDoList()
		if err != nil {
			fmt.Println(err)
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"todos": todos,
		})
	})

	e.POST("/todos", func(c *gin.Context) {
		title := c.PostForm("title")
		status := c.PostForm("status")

		id, err := CreateToDo(title, status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
			return
		}

		// return the new task as an HTML snippet so HTMX can update the page
		c.HTML(http.StatusOK, "task.html", gin.H{
			"Id":     id,
			"Title":  title,
			"Status": status,
		})
	})

	e.DELETE("/todos/:id", func(c *gin.Context) {
		param := c.Param("id")
		id, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = DeleteToDo(id)
		if err != nil {
			fmt.Println(err)
			return
		}

		// send 200 OK with no content (HTMX will remove the task)
		c.Status(http.StatusOK)
	})

	err := e.Run("localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("web server started at localhost:8080")

}
