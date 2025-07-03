package handler

import "github.com/gin-gonic/gin"

type TaskHandler struct{}

func (h *TaskHandler) RegisterRouter(router *gin.RouterGroup) {
	router.POST("/task", h.GenerateTask)
	router.GET("/task/:task_id", h.GetTaskStatus)
	router.PUT("/task/:task_id", h.ModifyTask)
}

func (h *TaskHandler) GenerateTask(c *gin.Context) {

}

func (h *TaskHandler) GetTaskStatus(c *gin.Context) {

}

func (h *TaskHandler) ModifyTask(c *gin.Context) {

}
