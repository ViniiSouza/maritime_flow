package leader

import (
	"encoding/json"
	"net/http"

	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/utils"
	"github.com/gin-gonic/gin"
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
	var request types.TowerHealthRequest
	if err := ctx.ShouldBindJSON(request); err != nil {
		utils.SetContextAndExecJSONWithErrorResponse(ctx, utils.ErrInvalidInput)
		return
	}

	if err := h.service.MarkTowerAsAlive(ctx, request.Id); err != nil {
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (h handler) AcquireSlot(ctx *gin.Context) {
	var request types.AcquireSlotRequest
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

func (h handler) ReleaseSlot(ctx *gin.Context) {
	var request types.ReleaseSlotLockRequest
	if err := ctx.ShouldBindJSON(request); err != nil {
		utils.SetContextAndExecJSONWithErrorResponse(ctx, utils.ErrInvalidInput)
		return
	}

	if err := h.service.ReleaseSlot(ctx, request); err != nil {
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
