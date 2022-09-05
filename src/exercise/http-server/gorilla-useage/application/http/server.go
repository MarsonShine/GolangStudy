package http

import (
	"encoding/json"
	"gorillaHttpServer/domain/taskstore"
	memorytaskstore "gorillaHttpServer/persistents/memory-taskstore"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type taskHttpServer struct {
	store taskstore.ITaskStore
}

func NewTaskHttpServer() *taskHttpServer {
	store := memorytaskstore.NewTaskStore()
	return &taskHttpServer{store: store}
}

func (ts *taskHttpServer) CreateTaskHandler(w http.ResponseWriter, req *http.Request) {

}

func (ts *taskHttpServer) GetAllTasksHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := ts.store.GetAllTasks()
	renderJson(w, allTasks)
}

func (ts *taskHttpServer) DeleteAllTasksHandler(w http.ResponseWriter, req *http.Request) {

}

func (ts *taskHttpServer) DeleteTaskHandler(w http.ResponseWriter, req *http.Request) {

}

func (ts *taskHttpServer) GetTaskHandler(w http.ResponseWriter, req *http.Request) {
	// gorilla自动解析路由
	log.Printf("handling get task at %s\n", req.URL.Path)
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	task := ts.store.GetTask(id)

	renderJson(w, task)
}

func renderJson[T any](w http.ResponseWriter, v T) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
