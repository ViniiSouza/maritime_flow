package minion

import (
	"log"
	"net/http"

	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/ViniiSouza/maritime_flow/com_tower/pkg/leaderelection"
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

func (h handler) ListTowers(ctx *gin.Context) {
	towers := h.service.ListTowers()
	response := types.TowersPayload{Towers: towers}

	ctx.JSON(http.StatusOK, response)
}

func (h handler) ListStructures(ctx *gin.Context) {
	structures := h.service.ListStructures()

	ctx.JSON(http.StatusOK, structures)
}

func (h handler) SyncTowers(ctx *gin.Context) {
	var towers types.TowersPayload
	if err := ctx.ShouldBindJSON(&towers); err != nil {
		log.Printf("failed to unmarshal request: %v", err)
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	h.service.SyncTowers(towers)
	ctx.JSON(http.StatusNoContent, nil)
}

func (h handler) SyncStructures(ctx *gin.Context) {
	var structures types.Structures
	if err := ctx.ShouldBindJSON(&structures); err != nil {
		log.Printf("failed to unmarshal request: %v", err)
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	h.service.SyncStructures(structures)
	ctx.JSON(http.StatusNoContent, nil)
}

func (h handler) CheckSlotAvailability(ctx *gin.Context) {
	var slotRequest types.SlotRequest
	if err := ctx.ShouldBindJSON(&slotRequest); err != nil {
		log.Printf("failed to unmarshal request: %v", err)
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	response, err := h.service.CheckSlotAvailability(ctx, slotRequest)
	if err != nil {
		log.Printf("failed to check slot availability: %v", err)
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h handler) HandleElection(ctx *gin.Context) {
	var req types.ElectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("failed to unmarshal request: %v", err)
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	uptime := config.Configuration.GetUptimeSeconds()

	var response types.ElectionResponse
	if uptime > req.CandidateUptime {
		log.Printf("[minion][election] my uptime (%.2fs) > candidate's uptime (%.2fs): starting my own election", uptime, req.CandidateUptime)
		if !config.Configuration.IsLeader() {
			go leaderelection.StartElection(h.service.ListTowers())
		}
		response = types.ElectionResponse{
			Uptime:          uptime,
			HasHigherUptime: true,
		}
	} else {
		log.Printf("my uptime (%.2fs) <= candidate's uptime (%.2fs): confirming vote in candidate", uptime, req.CandidateUptime)
		response = types.ElectionResponse{
			Uptime:          uptime,
			HasHigherUptime: false,
		}
	}

	ctx.JSON(http.StatusOK, response)
}

func (h handler) SetNewLeader(ctx *gin.Context) {
	var req types.NewLeaderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("failed to unmarshal request: %v", err)
		utils.SetContextAndExecJSONWithErrorResponse(ctx, err)
		return
	}

	config.Configuration.SetLeaderUUID(req.NewLeaderUUID)

	ctx.JSON(http.StatusNoContent, nil)
}
