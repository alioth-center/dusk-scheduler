package handler

import "github.com/gin-gonic/gin"

type TaskHandler struct{}

func (h *TaskHandler) RegisterRouter(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	v1.POST("/task", h.GenerateTask)
	v1.GET("/task/:task_id", h.GetTaskStatus)
	v1.PUT("/task/:task_id", h.ModifyTask)
}

func (h *TaskHandler) GenerateTask(c *gin.Context) {

}

func (h *TaskHandler) GetTaskStatus(c *gin.Context) {

}

func (h *TaskHandler) ModifyTask(c *gin.Context) {

}
