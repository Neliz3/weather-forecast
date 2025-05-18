package test

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func NewTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func NewTestServer(router *gin.Engine) *httptest.Server {
	return httptest.NewServer(router)
}
