package handler

import "github.com/gin-gonic/gin"

type ClientHandler struct{}

func (h *ClientHandler) RegisterHandler(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	v1.POST("/client", h.RegisterClient)
	v1.PUT("/client/:client_id/authorize", h.AuthorizeClient)
	v1.GET("/client/:client_id/metadata", h.GetMetadata)
	v1.GET("/client/:client_id/completed_tasks", h.GetCompletedTasks)
	v1.GET("/client/:client_id/quota_usage", h.GetQuotaUsage)
}

func (h *ClientHandler) RegisterClient(c *gin.Context) {}

func (h *ClientHandler) AuthorizeClient(c *gin.Context) {}

func (h *ClientHandler) GetMetadata(c *gin.Context) {}

func (h *ClientHandler) GetCompletedTasks(c *gin.Context) {}

func (h *ClientHandler) GetQuotaUsage(c *gin.Context) {}
