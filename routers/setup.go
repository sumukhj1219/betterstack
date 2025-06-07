package routers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sumukhj1219/betterstack/controllers"
)

func SetupGinRouter(ctx context.Context) *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.GET("/logs", func(c *gin.Context) {
		logs := controllers.GetMonitorLogs()
		c.JSON(http.StatusOK, logs)
	})

	return router
}
