package handler

import (
	"fmt"
	"github.com/alioth-center/dusk-scheduler/app/domain"
	"github.com/alioth-center/dusk-scheduler/app/entity"
	"github.com/alioth-center/dusk-scheduler/app/middleware"
	"github.com/alioth-center/dusk-scheduler/app/service"
	"github.com/alioth-center/dusk-scheduler/app/service/errors"
	"github.com/alioth-center/dusk-scheduler/infra/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	clientService  service.ClientService
	taskService    service.TaskService
	outcomeService service.OutcomeService
	brushService   service.BrushService
}

func NewTaskHandler(
	clientService service.ClientService,
	taskService service.TaskService,
	outcomeService service.OutcomeService,
	brushService service.BrushService,
) *TaskHandler {
	return &TaskHandler{
		clientService:  clientService,
		taskService:    taskService,
		outcomeService: outcomeService,
		brushService:   brushService,
	}
}

func (h *TaskHandler) RegisterHandler(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	v1.POST("/task", h.GenerateTask)
	v1.POST("/task/brush", h.GenerateBrushTask)
	v1.GET("/task/:task_id", h.GetTaskStatus)
}

func (h *TaskHandler) GenerateTask(c *gin.Context) {
	request := entity.CreateTaskRequest{}
	if bindErr := c.ShouldBindJSON(&request); bindErr != nil {
		errors.Ignore(c.Error(errors.BadRequestError(bindErr)))

		return
	}

	ctx, clientID := c.Request.Context(), c.GetUint64(middleware.CtxKeyClientID)
	client, existClient, queryClientErr := h.clientService.GetClientData(ctx, clientID)
	if queryClientErr != nil || !existClient {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}

	unifiedContent, unifyErr := utils.UnifyEncodingToBase64(request.Content, request.ContentEncoding)
	if unifyErr != nil {
		errors.Ignore(c.Error(errors.InvalidParameter("content")))

		return
	}

	task := domain.Task{
		Submitter:   clientID,
		ContentHash: utils.EncryptMd5String(unifiedContent),
		Width:       uint32(request.RenderWidth),
		Height:      uint32(request.RenderHeight),
		Type:        domain.TaskTypePainter,
		Priority:    client.AllowPriority,
		Format:      domain.TaskFormatFromString(request.OutputFormat),
		DelayRender: uint8(request.DelaySeconds),
	}
	taskID, createErr := h.taskService.CreateTask(ctx, &task, unifiedContent)
	if createErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}

	h.taskService.FlushSchedulerTaskQueue(ctx, taskID)
	c.Status(http.StatusCreated)
}

func (h *TaskHandler) GenerateBrushTask(c *gin.Context) {
	ctx, clientID := c.Request.Context(), c.GetUint64(middleware.CtxKeyClientID)
	widthStr, heightStr, delayStr := c.Query("width"), c.Query("height"), c.Query("delay")
	width, height, delay := uint32(0), uint32(0), uint8(0)
	if widthValue, convertErr := strconv.Atoi(widthStr); convertErr == nil {
		c.Status(http.StatusBadRequest)

		return
	} else {
		width = uint32(widthValue)
	}
	if heightValue, convertErr := strconv.Atoi(heightStr); convertErr == nil {
		c.Status(http.StatusBadRequest)

		return
	} else {
		height = uint32(heightValue)
	}
	if delayValue, convertErr := strconv.Atoi(delayStr); convertErr == nil {
		c.Status(http.StatusBadRequest)

		return
	} else {
		delay = uint8(delayValue)
	}

	client, existClient, queryClientErr := h.clientService.GetClientData(ctx, clientID)
	if queryClientErr != nil || !existClient {
		c.Status(http.StatusInternalServerError)

		return
	}

	contentData, readErr := c.GetRawData()
	if readErr != nil {
		c.Status(http.StatusBadRequest)

		return
	}

	encoded := utils.EncryptBase64String(string(contentData))
	task := domain.Task{
		Submitter:   clientID,
		ContentHash: utils.EncryptMd5String(encoded),
		Width:       width,
		Height:      height,
		Type:        domain.TaskTypeBrush,
		Priority:    client.AllowPriority,
		Format:      domain.TaskFormatRawImage,
		DelayRender: delay,
	}
	taskID, createErr := h.taskService.CreateTask(ctx, &task, encoded)
	if createErr != nil {
		c.Status(http.StatusInternalServerError)

		return
	}

	result, renderErr := h.brushService.RenderImage(ctx, taskID)
	if renderErr != nil {
		c.Status(http.StatusInternalServerError)

		return
	}

	c.Data(http.StatusOK, "image/png", result.Bytes())
}

func (h *TaskHandler) GetTaskStatus(c *gin.Context) {
	ctx, taskID := c.Request.Context(), uint64(0)
	if intVal, convertErr := strconv.Atoi(c.Param("task_id")); convertErr != nil {
		errors.Ignore(c.Error(errors.InvalidParameter("task_id")))

		return
	} else {
		taskID = uint64(intVal)
	}

	task, existTask, queryTaskErr := h.taskService.GetTaskByID(ctx, taskID)
	if queryTaskErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	} else if !existTask {
		errors.Ignore(c.Error(errors.GetTaskStatusNotFoundError()))

		return
	}

	response := entity.GetTaskStatusResponse{
		Size:        fmt.Sprintf("%dx%d", task.Width, task.Height),
		Priority:    task.Priority.String(),
		ContentHash: task.ContentHash,
		Status:      task.Status(),
		Timestamps: &entity.GetTaskStatusTimestamps{
			CreatedAt:   task.CreatedAt.Unix(),
			ScheduledAt: task.CreatedAt.Unix(),
			CompletedAt: task.CompletedAt.Unix(),
			ArchivedAt:  task.ArchivedAt.Unix(),
		},
	}
	switch task.Status() {
	case domain.TaskStatusCompleted:
		outcome, existOutcome, queryOutcomeErr := h.outcomeService.GetOutcomeByTaskID(ctx, task.ID)
		if queryOutcomeErr != nil || !existOutcome {
			errors.Ignore(c.Error(errors.InternalError()))

			return
		}

		response.OutcomeReference = outcome.Reference
	case domain.TaskStatusArchived:
		response.ArchiveReason = task.ArchiveReason.String()
	}

	c.JSON(http.StatusOK, entity.SuccessResponse(&response))
}
