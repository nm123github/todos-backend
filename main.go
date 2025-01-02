package main

import (
	"fx-todo-api/internal/redis"
	"fx-todo-api/internal/server"
	"fx-todo-api/internal/tasks"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(redis.NewRedisClient),
		fx.Provide(server.NewMux),
		fx.Provide(tasks.NewTaskHandler),
		fx.Invoke(server.StartServer),
	).Run()
}