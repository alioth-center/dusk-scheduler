package handler

import "github.com/gin-gonic/gin"

type Handler interface {
	RegisterHandler(router *gin.RouterGroup)
}
