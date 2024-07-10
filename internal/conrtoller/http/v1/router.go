package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	h := handler.Group("/v1")
	{
		h.GET("/heartbeat", heartBeat)
	}
}

func heartBeat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Heartbeat is ok!",
	})
}
