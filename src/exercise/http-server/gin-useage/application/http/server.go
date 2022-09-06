package http

import (
	"ginHttpServer/domain/taskstore"
	memorytaskstore "ginHttpServer/persistents/memory-taskstore"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type taskHttpServer struct {
	store taskstore.ITaskStore
}

func NewTaskHttpServer() *taskHttpServer {
	store := memorytaskstore.NewTaskStore()
	return &taskHttpServer{store: store}
}

func (ts *taskHttpServer) CreateTaskHandler(c *gin.Context) {
	type RequestTask struct {
		Text string    `json:"text"`
		Tags []string  `json:"tags"`
		Due  time.Time `json:"due"`
	}
	var rt RequestTask
	// from body/query bind parameters
	if err := c.ShouldBindJSON(&rt); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
	c.JSON(http.StatusOK, gin.H{"Id": id})
}

func (ts *taskHttpServer) GetAllTasksHandler(c *gin.Context) {
	allTasks := ts.store.GetAllTasks()
	// 框架自带json序列化
	c.JSON(http.StatusOK, allTasks)
}

func (ts *taskHttpServer) DeleteAllTasksHandler(c *gin.Context) {

}

func (ts *taskHttpServer) DeleteTaskHandler(c *gin.Context) {

}

func (ts *taskHttpServer) GetTaskHandler(c *gin.Context) {
	// gin web框架自动解析路由
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	task := ts.store.GetTask(id)
	c.JSON(http.StatusOK, task)
}
