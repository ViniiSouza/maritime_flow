package leader

import (
	"github.com/gin-gonic/gin"
)

func setupRouter(svc service) (router *gin.Engine) {
	handler := newHandler(svc)

	router = gin.Default()
	router.GET("towers", handler.ListHealthyTowers)
	router.GET("towers/", handler.ListHealthyTowers)
	router.POST("tower-health", handler.MarkTowerAsAlive)
	router.POST("tower-health/", handler.MarkTowerAsAlive)
	router.POST("acquire-slot", handler.AcquireSlot)
	router.POST("acquire-slot/", handler.AcquireSlot)
	router.POST("release-slot", handler.ReleaseSlot)
	router.POST("release-slot/", handler.ReleaseSlot)

	return
}
