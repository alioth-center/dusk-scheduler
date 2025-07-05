package handler

import "github.com/gin-gonic/gin"

type PainterHandler struct{}

func (h *PainterHandler) RegisterHandler(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	v1.POST("/painter", h.PainterConnect)
	v1.DELETE("/painter/:painter_id", h.PainterDisconnect)
	v1.GET("/painter/:painter_id/task", h.GetPaintTaskList)
	v1.POST("/painter/:painter_id/task/:task_id", h.CompleteTask)
}

func (h *PainterHandler) PainterConnect(c *gin.Context) {

}

func (h *PainterHandler) PainterDisconnect(c *gin.Context) {

}

func (h *PainterHandler) GetPaintTaskList(c *gin.Context) {

}

func (h *PainterHandler) CompleteTask(c *gin.Context) {

}
