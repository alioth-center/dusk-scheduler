package handler

import (
	"encoding/base64"
	"github.com/alioth-center/dusk-scheduler/app/domain"
	"github.com/alioth-center/dusk-scheduler/app/entity"
	"github.com/alioth-center/dusk-scheduler/app/service"
	"github.com/alioth-center/dusk-scheduler/app/service/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OutcomeHandler struct {
	taskService    service.TaskService
	outcomeService service.OutcomeService
	painterService service.PainterService
}

func (h *OutcomeHandler) RegisterHandler(router *gin.RouterGroup) {
	v1 := router.Group("/v1")
	v1.GET("/outcome/:outcome_reference", h.GetOutcomeContent)
	v1.PUT("/outcome/:outcome_reference/acknowledge", h.AcknowledgeOutcome)
}

func (h *OutcomeHandler) GetOutcomeContent(c *gin.Context) {
	outcomeReference := c.Param("outcome_reference")
	if len(outcomeReference) == 0 {
		errors.Ignore(c.Error(errors.InvalidParameter("outcome_reference")))

		return
	}

	ctx := c.Request.Context()
	outcome, existOutcome, getOutcomeErr := h.outcomeService.GetOutcomeContentByReference(ctx, outcomeReference)
	if getOutcomeErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}
	if !existOutcome {
		errors.Ignore(c.Error(errors.ResourceNotFound("outcome", outcomeReference)))

		return
	}

	task, existTask, getTaskErr := h.taskService.GetTaskByID(ctx, outcome.TaskID)
	if getTaskErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}
	if !existTask {
		errors.Ignore(c.Error(errors.ResourceNotFound("task", outcome.TaskID)))
	}

	painter, existPainter, getPainterErr := h.painterService.GetPainterByID(ctx, outcome.Instance)
	if getPainterErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}
	if !existPainter {
		errors.Ignore(c.Error(errors.ResourceNotFound("painter", outcome.Instance)))

		return
	}

	content := entity.GetOutcomeContent{}
	switch task.Format {
	case domain.TaskFormatImageURL:
		urlContent, getContentErr := h.outcomeService.GetOutcomeURL(ctx, outcome.Reference)
		if getContentErr != nil {
			errors.Ignore(c.Error(errors.InternalError()))

			return
		}

		content.ImageURL = urlContent.String()
	case domain.TaskFormatBase64Encoded:
		rawContent, getContentErr := h.outcomeService.GetOutcomeContent(ctx, outcomeReference)
		if getContentErr != nil {
			errors.Ignore(c.Error(errors.InternalError()))

			return
		}

		content.Base64Encoded = base64.StdEncoding.EncodeToString(rawContent.Bytes())
	case domain.TaskFormatRawImage:
		errors.Ignore(c.Error(errors.NotSupportRawOutputError()))

		return
	}

	response := entity.GetOutcomeResponse{
		Content: content,
		Metadata: entity.GetOutcomeMetadata{
			TaskReference: int(task.ID),
			InstanceName:  painter.Name,
			StartedAt:     outcome.StartedAt.Unix(),
			CompletedAt:   outcome.CompletedAt.Unix(),
		},
	}
	c.JSON(http.StatusOK, entity.SuccessResponse(&response))
}

func (h *OutcomeHandler) AcknowledgeOutcome(c *gin.Context) {
	outcomeReference := c.Param("outcome_reference")
	if len(outcomeReference) == 0 {
		errors.Ignore(c.Error(errors.InvalidParameter("outcome_reference")))

		return
	}

	ctx := c.Request.Context()
	existTask, archiveErr := h.taskService.ArchiveTaskByOutcomeReference(ctx, outcomeReference, domain.TaskArchiveReasonAcknowledged)
	if archiveErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}
	if !existTask {
		errors.Ignore(c.Error(errors.ResourceNotFound("outcome", outcomeReference)))

		return
	}

	c.Status(http.StatusNoContent)
}
