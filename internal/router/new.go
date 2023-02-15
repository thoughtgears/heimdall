package router

import (
	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	"github.com/thoughtgears/heimdall/handlers"
	"github.com/thoughtgears/heimdall/internal/config"
	"github.com/thoughtgears/heimdall/internal/gcp"
	"net/http"
)

func New(client *gcp.Client, config *config.Config) *gin.Engine {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(ginzerolog.Logger("gin"))

	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	v1 := r.Group("/v1")

	// Project endpoints
	v1.GET("/project", handlers.GetAll(client, config))

	return r
}
