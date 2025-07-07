package handler

import (
	"fmt"
	"github.com/alioth-center/dusk-scheduler/app/entity"
	"github.com/alioth-center/dusk-scheduler/app/middleware"
	"github.com/alioth-center/dusk-scheduler/app/service"
	"github.com/alioth-center/dusk-scheduler/app/service/errors"
	"github.com/alioth-center/dusk-scheduler/infra/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type ClientHandler struct {
	emailService       service.EmailService
	clientService      service.ClientService
	taskService        service.TaskService
	promotionalService service.PromotionalService
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
		errors.Ignore(c.Error(errors.BadRequestError(bindErr)))

		return
	}

	ctx := c.Request.Context()
	if validateErr := h.emailService.ValidateEmailAddress(ctx, request.EmailAddress); validateErr != nil {
		errors.Ignore(c.Error(errors.RegisterClientInvalidEmailAddressError()))

		return
	}

	existCode, queryCodeErr := h.promotionalService.CheckPromotionalByCode(ctx, request.RedemptionCode)
	if queryCodeErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}
	if !existCode {
		errors.Ignore(c.Error(errors.RedemptionCodeNotFoundError()))

		return
	}

	clientData, createErr := h.clientService.CreateClient(ctx, request.EmailAddress, request.RedemptionCode, c.ClientIP())
	if createErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}

	code := utils.GenerateAuthCode(6)
	expiredAt, cacheErr := h.clientService.StoreAuthorizationCode(ctx, clientData.ID, code)
	if cacheErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}

	args := map[string]any{"code": code}
	if sendErr := h.emailService.SendEmail(ctx, request.EmailAddress, service.EmailTemplateKeyRegisterClient, args); sendErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse(&entity.RegisterResponse{ExpiredAt: expiredAt.Unix()}))
}

func (h *ClientHandler) AuthorizeClient(c *gin.Context) {
	request := entity.AuthorizeRequest{}
	if bindErr := c.ShouldBindJSON(&request); bindErr != nil {
		errors.Ignore(c.Error(errors.BadRequestError(bindErr)))

		return
	}

	ctx, clientID := c.Request.Context(), uint64(0)
	if cid := c.Param("client_id"); len(cid) == 0 {
		errors.Ignore(c.Error(errors.BadRequestError(errors.InvalidParameter("client_id"))))

		return
	} else if intVal, convertErr := strconv.Atoi(cid); convertErr != nil {
		errors.Ignore(c.Error(errors.BadRequestError(errors.InvalidParameter("client_id"))))

		return
	} else {
		clientID = uint64(intVal)
	}

	authorized, maintainer, apiKey, authorizeErr := h.clientService.AuthorizeClient(ctx, clientID, request.EmailAddress, request.AuthorizationCode)
	if authorizeErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}
	if !authorized {
		errors.Ignore(c.Error(errors.AuthorizeClientFailedError()))

		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse(&entity.AuthorizeResponse{Maintainer: maintainer, ApiKey: apiKey}))
}

func (h *ClientHandler) GetMetadata(c *gin.Context) {
	ctx, clientID := c.Request.Context(), c.GetUint64(middleware.CtxKeyClientID)
	client, exist, queryErr := h.clientService.GetClientData(ctx, clientID)
	if queryErr != nil || !exist {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}

	response := entity.GetMetadataResponse{
		Maintainer: client.Maintainer,
		ApiKey:     client.ApiKey,
		Options: entity.GetMetadataClientOption{
			BrushApiEnable:          client.AllowBrush,
			DelayRenderEnable:       client.AllowDelay,
			MaxRenderHeight:         int(client.AllowHeight),
			MaxRenderWidth:          int(client.AllowWidth),
			MaxPriority:             int(client.AllowPriority),
			MaxRenderSize:           int(client.LimitRenderSize),
			RequestFrequency:        int(client.LimitFrequency),
			FrequencyIntervalSecond: int(client.LimitDuration),
		},
	}
	c.JSON(http.StatusOK, entity.SuccessResponse(&response))
}

func (h *ClientHandler) GetCompletedTasks(c *gin.Context) {
	filters, offsetTask := strings.Split(c.Query("filter"), ","), c.Query("offset_task")

	ctx, clientID, offsetTaskID := c.Request.Context(), c.GetUint64(middleware.CtxKeyClientID), uint64(0)
	if intVal, convertErr := strconv.Atoi(offsetTask); convertErr != nil {
		errors.Ignore(c.Error(errors.BadRequestError(errors.InvalidParameter("offset_task"))))

		return
	} else {
		offsetTaskID = uint64(intVal)
	}

	taskList, hasMore, queryErr := h.taskService.GetCompletedTasksByClientID(ctx, clientID, filters, offsetTaskID)
	if queryErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}

	response := entity.GetCompletedTasksResponse{
		HasMore: hasMore,
		Tasks:   make([]entity.GetCompletedTaskItem, len(taskList)),
	}
	for i, task := range taskList {
		response.Tasks[i] = entity.GetCompletedTaskItem{
			TaskID:        int(task.ID),
			Size:          fmt.Sprintf("%dx%d", int(task.Width), int(task.Height)),
			Priority:      int(task.Priority),
			ContentHash:   task.ContentHash,
			Status:        task.Status(),
			Timestamps:    entity.GetTaskStatusTimestamps{},
			ArchiveReason: task.ArchiveReason.String(),
		}
	}
	c.JSON(http.StatusOK, entity.SuccessResponse(&response))
}

func (h *ClientHandler) GetQuotaUsage(c *gin.Context) {
	ctx, clientID := c.Request.Context(), c.GetUint64(middleware.CtxKeyClientID)

	total, usage, checkTime, queryErr := h.clientService.GetClientQuotaUsage(ctx, clientID)
	if queryErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}

	response := entity.GetQuotaResponse{
		Details: entity.GetQuotaDetails{
			TotalQuota:     int(total),
			UsedQuota:      int(usage),
			RemainingQuota: int(total - usage),
		},
		LastCheckpoint: checkTime.Unix(),
	}
	c.JSON(http.StatusOK, entity.SuccessResponse(&response))
}
