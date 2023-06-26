package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// based on
type Todo struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Complete bool   `json:"complete"`
}

type TodoList struct {
	Todos []Todo `json:"todos"`
}

func newTodoList(initialTodos []Todo) *TodoList {
	var todoList TodoList = TodoList{
		Todos: append([]Todo{}, initialTodos...),
	}
	return &todoList
}

func (todoList *TodoList) getTodos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todoList.Todos)
}

func (todoList *TodoList) createTodo(c *gin.Context) {
	var newTodo Todo
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todoList.Todos = append(todoList.Todos, newTodo)
	c.JSON(http.StatusCreated, newTodo)
}

func (todoList *TodoList) updateTodoStatus(c *gin.Context) {
	todoID := c.Param("id")
	for i := range todoList.Todos {
		if todoList.Todos[i].ID == todoID {
			todoList.Todos[i].Complete = !todoList.Todos[i].Complete
			c.JSON(http.StatusOK, todoList.Todos)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
}

func (todoList *TodoList) deleteTodo(c *gin.Context) {
	todoID := c.Param("id")

	for i := range todoList.Todos {
		if todoList.Todos[i].ID == todoID {
			todoList.Todos = append(todoList.Todos[:i], todoList.Todos[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
}
