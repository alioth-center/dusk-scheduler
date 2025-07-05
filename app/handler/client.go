package handler

import (
	"github.com/alioth-center/dusk-scheduler/app/entity"
	"github.com/alioth-center/dusk-scheduler/app/service"
	"github.com/alioth-center/dusk-scheduler/app/service/errors"
	"github.com/alioth-center/dusk-scheduler/infra/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type ClientHandler struct {
	emailService    service.EmailService
	locationService service.LocationService
	clientService   service.ClientService
}

func (h *ClientHandler) RegisterHandler(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	v1.POST("/client", h.RegisterClient)
	v1.PUT("/client/:client_id/authorize", h.AuthorizeClient)
	v1.GET("/client/:client_id/metadata", h.GetMetadata)
	v1.GET("/client/:client_id/completed_tasks", h.GetCompletedTasks)
	v1.GET("/client/:client_id/quota_usage", h.GetQuotaUsage)
}

func (h *ClientHandler) RegisterClient(c *gin.Context) {
	request := entity.RegisterRequest{}
	if bindErr := c.ShouldBindJSON(&request); bindErr != nil {
		_ = c.Error(errors.BadRequestError(bindErr))

		return
	}

	ctx := c.Request.Context()
	if validateErr := h.emailService.ValidateEmailAddress(ctx, request.EmailAddress); validateErr != nil {
		_ = c.Error(errors.RegisterClientInvalidEmailAddressError())

		return
	}

	clientData, createErr := h.clientService.CreateClient(ctx, request.EmailAddress, request.RedemptionCode, c.ClientIP())
	if createErr != nil {
		_ = c.Error(errors.InternalError())

		return
	}

	code := strings.ToUpper(utils.GenerateRandomString(6))
	expiredAt, cacheErr := h.clientService.StoreAuthorizationCode(ctx, clientData.ID, code)
	if cacheErr != nil {
		_ = c.Error(errors.InternalError())

		return
	}

	args := map[string]any{"code": code}
	if sendErr := h.emailService.SendEmail(ctx, request.EmailAddress, service.EmailTemplateKeyRegisterClient, args); sendErr != nil {
		_ = c.Error(errors.InternalError())

		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse(&entity.RegisterResponse{ExpiredAt: expiredAt.Unix()}))
}

func (h *ClientHandler) AuthorizeClient(c *gin.Context) {}

func (h *ClientHandler) GetMetadata(c *gin.Context) {}

func (h *ClientHandler) GetCompletedTasks(c *gin.Context) {}

func (h *ClientHandler) GetQuotaUsage(c *gin.Context) {}
