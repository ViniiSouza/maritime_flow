package minion

import (
	"encoding/json"
	"net/http"

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

	response, err := json.Marshal(tower.TowersResponse{Towers: towers})
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
