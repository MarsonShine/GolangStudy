package taskstore

import "time"

// refference by https://eli.thegreenplace.net/2021/rest-servers-in-go-part-1-standard-library/
type Task struct {
	Id          int
	Text        string
	Due         time.Time
	Tags        []string
	Attachments []Attachment
}
type Attachment struct {
	Name     string    `json:"Name"`
	Date     time.Time `json:"Date"`
	Contents string    `json:"Contents"`
}

type ITaskStore interface {
	CreateTask(txt string, tags []string, due time.Time, attachments []Attachment) int
	DeleteTask(id int) error
	GetAllTasks() []Task
	GetTask(id int) Task
	GetTasksByTag(tag string) []Task
	GetTasksByDueDate(year int, month time.Month, day int) []Task
}
