package tasks

import (
	"context"
	"fmt"
	"fx-todo-api/internal/redis"
	"fx-todo-api/pkg"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type TaskHandler struct {
	Redis *redis.RedisClient
}

func NewTaskHandler(redis *redis.RedisClient) *TaskHandler {
	return &TaskHandler{Redis: redis}
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/task/"):]
	if id == "" {
		http.Error(w, "Missing task ID", http.StatusBadRequest)
		return
	}

	if err := h.Redis.Client.HSet(context.Background(), fmt.Sprintf("task:%s", id), "status", "completed").Err(); err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	pkg.WriteJSON(w, map[string]string{"message": "Task marked as complete"})
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/task/"):]
	if id == "" {
		http.Error(w, "Missing task ID", http.StatusBadRequest)
		return
	}

	if err := h.Redis.Client.Del(context.Background(), fmt.Sprintf("task:%s", id)).Err(); err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	pkg.WriteJSON(w, map[string]string{"message": "Task deleted"})
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := pkg.ParseJSON(r, &task); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	task.ID = uuid.New().String()
	task.DateCreated = time.Now().Format(time.RFC3339)
	task.Status = "in_progress"

	key := fmt.Sprintf("task:%s", task.ID)
	if err := h.Redis.Client.HSet(context.Background(), key, map[string]interface{}{
		"id":        task.ID,
		"name":        task.Name,
		"status":      task.Status,
		"dateCreated": task.DateCreated,
	}).Err(); err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	pkg.WriteJSON(w, task)
}

func (h *TaskHandler) ListTask(w http.ResponseWriter, r *http.Request) {
	keys, err := h.Redis.Client.Keys(context.Background(), "task:*").Result()
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	var tasks []Task
	for _, key := range keys {
		data, err := h.Redis.Client.HGetAll(context.Background(), key).Result()
		if err != nil {
			continue
		}

		tasks = append(tasks, Task{
			ID:        data["id"],
			Name:      data["name"],
			Status:    data["status"],
			DateCreated: data["dateCreated"],
		})
	}

	fmt.Println("Tasks > ", tasks)

	pkg.WriteJSON(w, tasks)
}
