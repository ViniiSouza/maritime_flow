package leader

import (
	"github.com/gin-gonic/gin"
)

func setupRouter() (router *gin.Engine) {
	repo := newRepository()
	svc := newService(repo)
	handler := newHandler(svc)

	router = gin.Default()
	router.POST("tower-health", handler.MarkTowerAsAlive)
	router.POST("acquire-slot", handler.AcquireSlot)
	router.POST("release-slot", handler.ReleaseSlot)

	return
}
