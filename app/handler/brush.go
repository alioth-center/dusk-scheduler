package handler

import (
	"github.com/alioth-center/dusk-scheduler/app/entity"
	"github.com/alioth-center/dusk-scheduler/app/service"
	"github.com/alioth-center/dusk-scheduler/app/service/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	v1.DELETE("/brush/:brush_id", h.BrushDisconnect)
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

	brushID, createErr := h.brushService.CreateBrush(ctx, request.Maintainer, request.Protocol, request.CallURL)
	if createErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}

	c.JSON(http.StatusOK, entity.SuccessResponse(&entity.RegisterBrushResponse{BrushID: int(brushID)}))
}

func (h *BrushHandler) BrushDisconnect(c *gin.Context) {
	ctx, brushID := c.Request.Context(), uint64(0)
	if params := c.Param("brush_id"); len(params) == 0 {
		errors.Ignore(c.Error(errors.InvalidParameter("brush_id")))

		return
	} else if intVal, convertErr := strconv.ParseUint(params, 10, 64); convertErr != nil {
		errors.Ignore(c.Error(errors.InvalidParameter("brush_id")))

		return
	} else {
		brushID = intVal
	}

	if disconnectErr := h.brushService.DisconnectBrush(ctx, brushID); disconnectErr != nil {
		errors.Ignore(c.Error(errors.InternalError()))

		return
	}

	c.Status(http.StatusNoContent)
}
