package leader

import (
	"encoding/json"
	"net/http"

	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/slot"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/tower"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handler struct {
	service service
}

func newHandler(service service) handler {
	return handler{
		service: service,
	}
}

func (h handler) MarkTowerAsAlive(ctx *gin.Context) {
	var request tower.TowerHealthRequest
	if err := ctx.ShouldBindJSON(request); err != nil {
		utils.SetContextAndExecJSONWithErrorResponse(ctx, utils.ErrInvalidInput)
		return
	}

	id, err := uuid.Parse(request.Id)
	if err != nil {
		utils.SetContextAndExecJSONWithErrorResponse(ctx, utils.ErrInvalidUUID)
		return
	}

	if err := h.service.MarkTowerAsAlive(ctx, id); err != nil {
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (h handler) AcquireSlot(ctx *gin.Context) {
	var request slot.AcquireSlotRequest
	if err := ctx.ShouldBindJSON(request); err != nil {
		utils.SetContextAndExecJSONWithErrorResponse(ctx, utils.ErrInvalidInput)
		return
	}

	result, err := h.service.AcquireSlot(ctx, request)
	if err != nil {
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	response, err := json.Marshal(result)
	if err != nil {
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
