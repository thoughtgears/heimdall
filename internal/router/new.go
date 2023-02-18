package router

import (
	"net/http"

	"github.com/thoughtgears/heimdall/handlers/service_cataloge"

	"cloud.google.com/go/firestore"
	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	"github.com/thoughtgears/heimdall/handlers/project"
	"github.com/thoughtgears/heimdall/internal/config"
)

func New(client *firestore.Client, config *config.Config) *gin.Engine {
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
	v1.GET("/projects", project.GetAll(client, config))
	v1.GET("/projects/:id", project.GetID(client, config))

	// Service Catalogue endpoints
	v1.GET("/catalogue", service_cataloge.Get(client))
	v1.POST("/catalogue", service_cataloge.Post(client))

	return r
}
