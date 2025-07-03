package handler

import "github.com/gin-gonic/gin"

type OutcomeHandler struct{}

func (h *OutcomeHandler) RegisterHandler(router *gin.RouterGroup) {
	router.GET("/outcome/:outcome_id", h.GetOutcomeContent)
	router.PUT("/outcome/:outcome_id/acknowledge", h.AcknowledgeOutcome)
}

func (h *OutcomeHandler) GetOutcomeContent(c *gin.Context) {

}

func (h *OutcomeHandler) AcknowledgeOutcome(c *gin.Context) {

}
