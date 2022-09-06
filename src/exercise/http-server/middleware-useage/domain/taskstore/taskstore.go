package taskstore

import "time"

// refference by https://eli.thegreenplace.net/2021/rest-servers-in-go-part-1-standard-library/
type Task struct {
	Id   int
	Text string
	Due  time.Time
	Tags []string
}

type ITaskStore interface {
	CreateTask(txt string, tags []string, due time.Time) int
	DeleteTask(id int) error
	GetAllTasks() []Task
	GetTasksByTag(tag string) []Task
	GetTasksByDueDate(year int, month time.Month, day int) []Task
}
