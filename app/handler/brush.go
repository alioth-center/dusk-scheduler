package handler

import (
	"github.com/alioth-center/dusk-scheduler/app/entity"
	"github.com/alioth-center/dusk-scheduler/app/service"
	"github.com/alioth-center/dusk-scheduler/app/service/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BrushHandler struct {
	emailService service.EmailService
	brushService service.BrushService
}

func NewBrushHandler(
	emailService service.EmailService,
	brushService service.BrushService,
) *BrushHandler {
	return &BrushHandler{
		emailService: emailService,
		brushService: brushService,
	}
}

func (h *BrushHandler) RegisterHandler(router *gin.RouterGroup) {
	v1 := router.Group("v1")
	v1.POST("/brush", h.BrushConnect)
	v1.DELETE("/brush/:brush_name", h.BrushDisconnect)
}

func (h *BrushHandler) BrushConnect(c *gin.Context) {
	request := entity.RegisterBrushRequest{}
	if bindErr := c.ShouldBind(&request); bindErr != nil {
		errors.Ignore(c.Error(errors.BadRequestError(bindErr)))

		return
	}

	ctx := c.Request.Context()
	if validateErr := h.emailService.ValidateEmailAddress(ctx, request.Maintainer); validateErr != nil {
		errors.Ignore(c.Error(errors.RegisterPainterInvalidEmailAddressError()))

		return
	}

	painter, policy, createErr := h.brushService.CreateBrush(ctx, request.Maintainer, request.Protocol)
	if createErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}

	response := entity.RegisterBrushResponse{
		Name:   painter.Name,
		Secret: painter.Secret,
		Policy: entity.RegisterBrushPolicy{
			Protocol: policy.Protocol.String(),
			Options:  policy.Options,
		},
	}
	c.JSON(http.StatusOK, entity.SuccessResponse(&response))
}

func (h *BrushHandler) BrushDisconnect(c *gin.Context) {
	ctx, brushName := c.Request.Context(), c.Param("brush_name")
	if disconnectErr := h.brushService.DisconnectBrush(ctx, brushName); disconnectErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}

	c.Status(http.StatusNoContent)
}
