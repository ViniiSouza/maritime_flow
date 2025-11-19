package minion

import (
	"github.com/gin-gonic/gin"
)

func setupRouter() (router *gin.Engine) {
	repo := newRepository()
	svc := newService(repo)
	handler := newHandler(svc)

	router = gin.Default()
	router.GET("/towers", handler.ListTowers)
	router.GET("/structures", handler.ListStructures)

	return
}
