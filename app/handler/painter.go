package handler

import (
	"github.com/alioth-center/dusk-scheduler/app/domain"
	"github.com/alioth-center/dusk-scheduler/app/entity"
	"github.com/alioth-center/dusk-scheduler/app/middleware"
	"github.com/alioth-center/dusk-scheduler/app/service"
	"github.com/alioth-center/dusk-scheduler/app/service/errors"
	"github.com/alioth-center/dusk-scheduler/infra/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type PainterHandler struct {
	taskService    service.TaskService
	outcomeService service.OutcomeService
	emailService   service.EmailService
	painterService service.PainterService
}

func NewPainterHandler(
	taskService service.TaskService,
	outcomeService service.OutcomeService,
	emailService service.EmailService,
	painterService service.PainterService,
) *PainterHandler {
	return &PainterHandler{
		taskService:    taskService,
		outcomeService: outcomeService,
		emailService:   emailService,
		painterService: painterService,
	}
}

func (h *PainterHandler) RegisterHandler(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	v1.POST("/painter", h.PainterConnect)
	v1.PUT("/painter/:painter_name", h.PainterReconnect)
	v1.DELETE("/painter/:painter_name", h.PainterDisconnect)
	v1.GET("/painter/:painter_name/task", h.GetPaintTaskList)
	v1.POST("/painter/:painter_name/task/:task_id", h.CompleteTask)
}

func (h *PainterHandler) PainterConnect(c *gin.Context) {
	request := entity.RegisterPainterRequest{}
	if bindErr := c.ShouldBindJSON(&request); bindErr != nil {
		errors.Ignore(c.Error(errors.BadRequestError(bindErr)))

		return
	}

	ctx := c.Request.Context()
	if validateErr := h.emailService.ValidateEmailAddress(ctx, request.Maintainer); validateErr != nil {
		errors.Ignore(c.Error(errors.RegisterPainterInvalidEmailAddressError()))

		return
	}

	painter, policy, createErr := h.painterService.CreatePainter(ctx, request.Maintainer, request.Slot, c.ClientIP())
	if createErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}
	h.taskService.FlushPainterScheduler(ctx, painter.Name)

	response := entity.RegisterPainterResponse{
		Name:   painter.Name,
		Secret: painter.Secret,
		Policy: entity.RegisterPainterPolicy{
			Protocol: policy.Protocol.String(),
			Options:  policy.Options,
		},
	}
	c.JSON(http.StatusOK, entity.SuccessResponse(&response))
}

func (h *PainterHandler) PainterReconnect(c *gin.Context) {
	ctx, painterName := c.Request.Context(), c.Param("painter_name")
	isHeartbeat, connectErr := h.painterService.ReconnectPainter(ctx, painterName)
	if connectErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}

	if !isHeartbeat {
		h.taskService.FlushPainterScheduler(ctx, painterName)
	}

	c.Status(http.StatusNoContent)
}

func (h *PainterHandler) PainterDisconnect(c *gin.Context) {
	ctx, painterName := c.Request.Context(), c.Param("painter_name")
	if disconnectErr := h.painterService.DisconnectPainter(ctx, painterName); disconnectErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}

	h.taskService.FlushPainterScheduler(ctx, painterName)
	c.Status(http.StatusNoContent)
}

func (h *PainterHandler) GetPaintTaskList(c *gin.Context) {
	ctx, painterName, secret := c.Request.Context(), c.Param("painter_name"), c.GetString(middleware.CtxKeyPainterSecret)
	taskList, taskContent, getTaskErr := h.taskService.GetScheduledTaskListByPainterName(ctx, painterName)
	if getTaskErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}

	tasks := make([]entity.GetPainterTaskListItem, len(taskList))
	for i, task := range taskList {
		content := taskContent[i]
		tasks[i] = entity.GetPainterTaskListItem{
			TaskID:         int(task.ID),
			Height:         int(task.Height),
			Width:          int(task.Width),
			Priority:       int(task.Priority),
			DelaySeconds:   int(task.DelayRender),
			EncodedContent: content,
			Checksum:       utils.EncryptHmacSha256String(content, secret),
		}
	}
	c.JSON(http.StatusOK, entity.SuccessResponse(&entity.GetPainterTaskListResponse{Tasks: tasks}))
}

func (h *PainterHandler) CompleteTask(c *gin.Context) {
	request := entity.CompletePainterTaskRequest{}
	if bindErr := c.ShouldBindJSON(&request); bindErr != nil {
		errors.Ignore(c.Error(errors.BadRequestError(bindErr)))

		return
	}

	ctx, painterName, painterSecret, task := c.Request.Context(), c.Param("painter_name"), c.GetString(middleware.CtxKeyPainterSecret), &domain.Task{}
	if intVal, convertErr := strconv.ParseUint(c.Param("task_id"), 10, 64); convertErr != nil {
		errors.Ignore(c.Error(errors.InvalidParameter("task_id")))

		return
	} else if checkedTask, exist, checkErr := h.taskService.GetTaskByID(ctx, intVal); checkErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	} else if !exist {
		errors.Ignore(c.Error(errors.CompleteTaskNotFoundTaskError()))

		return
	} else {
		task = checkedTask
	}

	if completeErr := h.taskService.CompleteTask(ctx, task.ID); completeErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}

	if request.Status == domain.OutcomeCompleteReasonCompleted.String() {
		startedAt, completedAt := time.Unix(request.StartedAt, 0), time.Unix(request.CompletedAt, 0)
		if createErr := h.outcomeService.CreateOutcome(ctx, painterName, task.ID, request.StorageReference, startedAt, completedAt); createErr != nil {
			errors.Ignore(c.Error(errors.InternalError()))

			return
		}
	}

	nextList, nextContent, queryErr := h.taskService.GetNextScheduledTaskListByPainterName(ctx, painterName)
	if queryErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}

	taskList := make([]entity.GetPainterTaskListItem, len(nextList))
	for i, taskItem := range nextList {
		content := nextContent[i]
		taskList[i] = entity.GetPainterTaskListItem{
			TaskID:         int(taskItem.ID),
			Height:         int(taskItem.Height),
			Width:          int(taskItem.Width),
			Priority:       int(taskItem.Priority),
			DelaySeconds:   int(taskItem.DelayRender),
			EncodedContent: content,
			Checksum:       utils.EncryptHmacSha256String(content, painterSecret),
		}
	}
	c.JSON(http.StatusOK, entity.SuccessResponse(&entity.GetPainterTaskListResponse{Tasks: taskList}))
}
