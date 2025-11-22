package types

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

type BodyIntercepter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w BodyIntercepter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}
