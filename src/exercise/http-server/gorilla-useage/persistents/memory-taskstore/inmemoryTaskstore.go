package memorytaskstore

import (
	"fmt"
	"gorillaHttpServer/domain/taskstore"
	"sync"
	"time"
)

type InMemoryTaskStore struct {
	sync.Mutex
	tasks  map[int]taskstore.Task
	nextId int
}

// GetTask implements taskstore.ITaskStore
func (*InMemoryTaskStore) GetTask(id int) taskstore.Task {
	panic("unimplemented")
}

// CreateTask implements taskstore.ITaskStore
func (ts *InMemoryTaskStore) CreateTask(text string, tags []string, due time.Time) int {
	ts.Lock()
	defer ts.Unlock()
	task := taskstore.Task{
		Id:   ts.nextId,
		Text: text,
		Due:  due,
	}
	task.Tags = make([]string, len(tags))
	copy(task.Tags, tags)

	ts.tasks[ts.nextId] = task
	ts.nextId++
	return task.Id
}

// DeleteTask implements taskstore.ITaskStore
func (ts *InMemoryTaskStore) DeleteTask(id int) error {
	ts.Lock()
	defer ts.Unlock()
	if _, ok := ts.tasks[id]; !ok {
		return fmt.Errorf("task with id=%d not found", id)
	}
	delete(ts.tasks, id)
	return nil
}

// GetAllTasks implements taskstore.ITaskStore
func (ts *InMemoryTaskStore) GetAllTasks() []taskstore.Task {
	tasks := make([]taskstore.Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// GetTasksByDueDate implements taskstore.ITaskStore
func (*InMemoryTaskStore) GetTasksByDueDate(year int, month time.Month, day int) []taskstore.Task {
	panic("unimplemented")
}

// GetTasksByTag implements taskstore.ITaskStore
func (*InMemoryTaskStore) GetTasksByTag(tag string) []taskstore.Task {
	panic("unimplemented")
}

func NewTaskStore() taskstore.ITaskStore {
	ts := InMemoryTaskStore{}
	ts.tasks = make(map[int]taskstore.Task)
	ts.nextId = 0
	return &ts
}
