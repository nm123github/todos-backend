package tasks

type Task struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status     string `json:"status"` // in_progress, completed
	DateCreated string `json:"dateCreated"`
}
