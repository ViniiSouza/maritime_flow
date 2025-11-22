package leader

import (
	"github.com/gin-gonic/gin"
)

func setupRouter(svc service) (router *gin.Engine) {
	handler := newHandler(svc)

	router = gin.Default()
	router.POST("tower-health", handler.MarkTowerAsAlive)
	router.POST("acquire-slot", handler.AcquireSlot)
	router.POST("release-slot", handler.ReleaseSlot)

	return
}
