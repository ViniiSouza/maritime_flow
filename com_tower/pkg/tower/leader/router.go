package leader

import (
	"github.com/ViniiSouza/maritime_flow/com_tower/config"
	"github.com/gin-gonic/gin"
)

func setupRouter(cfg config.Config) (router *gin.Engine) {
	repo := newRepository(cfg.DB)
	svc := newService(repo)
	handler := newHandler(svc)

	router = gin.Default()
	router.POST("tower-health", handler.MarkTowerAsAlive)

	return
}
