package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	h := handler.Group("/v1")
	{
		h.GET("/heartbeat", heartBeat)
		h.GET("/greater/:name", greater)
	}
}

func heartBeat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Heartbeat is ok!",
	})
}

func greater(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Hello, %s", c.Param("name")),
	})

}
