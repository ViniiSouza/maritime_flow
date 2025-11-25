package leader

import (
	"log"
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
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("failed to unmarshal request: %v", err)
		utils.SetContextAndExecJSONWithErrorResponse(ctx, utils.ErrInvalidInput)
		return
	}

	if err := h.service.MarkTowerAsAlive(ctx, request.Id); err != nil {
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (h handler) ListHealthyTowers(ctx *gin.Context) {
	towers, err := h.service.ListHealthyTowers(ctx)
	if err != nil {
		log.Printf("failed to list healthy towers: %v", err)
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	response := types.TowersPayload{Towers: towers}
	ctx.JSON(http.StatusOK, response)
}

func (h handler) AcquireSlot(ctx *gin.Context) {
	var request types.AcquireSlotRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("failed to unmarshal request: %v", err)
		utils.SetContextAndExecJSONWithErrorResponse(ctx, utils.ErrInvalidInput)
		return
	}

	response, err := h.service.AcquireSlot(ctx, request)
	if err != nil {
		log.Printf("failed to acquire slot: %v", err)
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h handler) ReleaseSlot(ctx *gin.Context) {
	var request types.ReleaseSlotLockRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("failed to unmarshal request: %v", err)
		utils.SetContextAndExecJSONWithErrorResponse(ctx, utils.ErrInvalidInput)
		return
	}

	if err := h.service.ReleaseSlot(ctx, request); err != nil {
		log.Printf("failed to release slot: %v", err)
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
