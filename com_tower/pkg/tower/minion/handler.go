package minion

import (
	"encoding/json"
	"net/http"

	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/types"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/tower"
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

func (h handler) ListTowers(ctx *gin.Context) {
	towers := h.service.ListTowers()

	response, err := json.Marshal(tower.TowersPayload{Towers: towers})
	if err != nil {
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h handler) ListStructures(ctx *gin.Context) {
	structures := h.service.ListStructures()

	response, err := json.Marshal(structures)
	if err != nil {
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h handler) SyncTowers(ctx *gin.Context) {
	var towers tower.TowersPayload
	if err := ctx.ShouldBindJSON(towers); err != nil {
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	h.service.SyncTowers(towers)
	ctx.JSON(http.StatusNoContent, nil)
}

func (h handler) SyncStructures(ctx *gin.Context) {
	var structures types.Structures
	if err := ctx.ShouldBindJSON(structures); err != nil {
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	h.service.SyncStructures(structures)
	ctx.JSON(http.StatusNoContent, nil)
}

func (h handler) CheckSlotAvailability(ctx *gin.Context) {
	var slotRequest types.SlotRequest
	if err := ctx.ShouldBindJSON(slotRequest); err != nil {
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	result, err := h.service.CheckSlotAvailability(ctx, slotRequest)
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
