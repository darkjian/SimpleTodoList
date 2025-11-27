package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/darkjian/simpletodolist/internal/adapters/database"
	"github.com/darkjian/simpletodolist/internal/adapters/delivery"
	"github.com/darkjian/simpletodolist/internal/config"
	"github.com/darkjian/simpletodolist/internal/usecase"
	"github.com/gin-gonic/gin"
)

func Run() {

	conf, err := config.Load("./config/local.yaml")
	if err != nil {
		slog.Error("loading config", slog.Any("error", err))
		os.Exit(1)
	}

	db, err := database.NewPostgresConnection(conf.DatabaseUrl)
	if err != nil {
		slog.Error("postgres init", slog.Any("error", err))
		os.Exit(1)
	}

	task_handler := delivery.NewTaskHandler(
		usecase.NewTaskService(
			database.NewTaskRepository(db),
		),
	)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// api
	api := r.Group("/api/v1")
	{
		api.GET("/tasks", task_handler.ListTasks)
		api.GET("/tasks/:id", task_handler.GetTask)
		api.POST("/tasks", task_handler.CreateTask)
		// api.PUT("/tasks/:id", task_handler.UpdateTask)
		api.DELETE("/tasks/:id", task_handler.DeleteTask)
	}

	server := &http.Server{
		Addr:         conf.Server.Port,
		Handler:      r,
		ReadTimeout:  time.Duration(conf.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(conf.Server.IdleTimeout) * time.Second,
	}

	slog.Info("starting server", slog.String("port", conf.Server.Port))
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("starting server", slog.Any("error", err))
		}
	}()

	ch := make(chan os.Signal, 1)
	defer close(ch)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("shutdown server", slog.Any("error", err))
	}

	slog.Info("shutdown server gracefully")
}
