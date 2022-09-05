package http

import (
	"basicHttpServer/domain/taskstore"
	memorytaskstore "basicHttpServer/persistents/memory-taskstore"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type taskHttpServer struct {
	store *taskstore.ITaskStore
}

func NewTaskHttpServer() *taskHttpServer {
	store := memorytaskstore.NewTaskStore()
	return &taskHttpServer{store: &store}
}

func (ts *taskHttpServer) TaskHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/task/" {
		if req.Method == http.MethodPost {
			ts.createTaskHandler(w, req)
		} else if req.Method == http.MethodGet {
			ts.getAllTasksHandler(w, req)
		} else if req.Method == http.MethodDelete {
			ts.deleteAllTasksHandler(w, req)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, DELETE or POST at /task/, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {
		// Request has an ID, as in "/task/<id>".
		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			http.Error(w, "expect /task/<id> in task handler", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Method == http.MethodDelete {
			ts.deleteTaskHandler(w, req, int(id))
		} else if req.Method == http.MethodGet {
			ts.getTaskHandler(w, req, int(id))
		} else {
			http.Error(w, fmt.Sprintf("expect method GET or DELETE at /task/<id>, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}

func (ts *taskHttpServer) createTaskHandler(w http.ResponseWriter, req *http.Request) {

}

func (ts *taskHttpServer) getAllTasksHandler(w http.ResponseWriter, req *http.Request) {

}

func (ts *taskHttpServer) deleteAllTasksHandler(w http.ResponseWriter, req *http.Request) {

}

func (ts *taskHttpServer) deleteTaskHandler(w http.ResponseWriter, req *http.Request, taskId int) {

}

func (ts *taskHttpServer) getTaskHandler(w http.ResponseWriter, req *http.Request, taskId int) {

}
