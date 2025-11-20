package minion

import (
	"github.com/gin-gonic/gin"
)

func setupRouter() (router *gin.Engine) {
	repo := newRepository()
	integ := newIntegration()
	svc := newService(integ, repo)
	handler := newHandler(svc)

	router = gin.Default()
	router.GET("/towers", handler.ListTowers)
	router.GET("/structures", handler.ListStructures)
	router.POST("/towers", handler.SyncTowers)
	router.POST("/structures", handler.SyncStructures)
	router.POST("/slots", handler.CheckSlotAvailability)

	return
}
