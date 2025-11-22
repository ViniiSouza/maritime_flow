package minion

import (
	"github.com/gin-gonic/gin"
)

func setupRouter(svc service) (router *gin.Engine) {
	handler := newHandler(svc)

	router = gin.Default()
	router.Use(AuditRequests())
	router.GET("/towers", handler.ListTowers)
	router.GET("/structures", handler.ListStructures)
	router.POST("/towers", handler.SyncTowers)
	router.POST("/structures", handler.SyncStructures)
	router.POST("/slots", handler.CheckSlotAvailability)
	router.POST("/election", handler.HandleElection)
	router.POST("/leader", handler.SetNewLeader)

	return
}
