package handler

import "github.com/gin-gonic/gin"

type PainterHandler struct{}

func (h *PainterHandler) RegisterHandler(router *gin.RouterGroup) {
	router.POST("/painter", h.PainterConnect)
	router.DELETE("/painter/:painter_id", h.PainterDisconnect)
	router.GET("/painter/:painter_id/task", h.GetPaintTaskList)
	router.POST("/painter/:painter_id/task/:task_id", h.CompleteTask)
}

func (h *PainterHandler) PainterConnect(c *gin.Context) {

}

func (h *PainterHandler) PainterDisconnect(c *gin.Context) {

}

func (h *PainterHandler) GetPaintTaskList(c *gin.Context) {

}

func (h *PainterHandler) CompleteTask(c *gin.Context) {

}
