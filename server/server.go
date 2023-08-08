package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/swashbuck1r/simple-go-api-server/config"
	"github.com/swashbuck1r/simple-go-api-server/server/middleware"
)

var (
	SignalChan chan os.Signal
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func Start(c *config.Configuration) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Handle ctrl-C signals
	SignalChan = make(chan os.Signal, 1)
	signal.Notify(SignalChan, os.Interrupt)
	defer func() {
		signal.Stop(SignalChan)
		cancel()
	}()
	go func() {
		select {
		case <-SignalChan: // first signal, cancel context
			log.Printf("Received an interrupt, stopping services...")
			cancel()
		case <-ctx.Done():
		}
		os.Exit(2)
	}()

	r := setupRouter(ctx)
	log.Printf("Server running on port %d", c.Port)
	if err := r.Run(fmt.Sprintf(":%d", c.Port)); err != nil {
		log.Fatal("could not run server", err)
	}
}

func setupRouter(ctx context.Context) *gin.Engine {
	ginEngine := gin.New()

	ginEngine.Use(gin.Recovery())
	ginEngine.Use(middleware.HTTPLogger())

	initialTodos := []Todo{}
	viper.UnmarshalKey("InitialToDos", &initialTodos)

	todoList := newTodoList(initialTodos)
	ginEngine.GET("/", todoList.getTodos)
	ginEngine.GET("/todos", todoList.getTodos)
	ginEngine.POST("/todos", todoList.createTodo)
	ginEngine.PUT("/todos/:id", todoList.updateTodoStatus)
	ginEngine.DELETE("/todos/:id", todoList.deleteTodo)

	return ginEngine
}
