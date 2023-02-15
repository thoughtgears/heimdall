package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/thoughtgears/heimdall/internal/config"
	"github.com/thoughtgears/heimdall/internal/gcp"
	"net/http"
)

func GetAll(client *gcp.Client, config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := client.GetAllProjects(c, config.Collection)
		if err != nil {
			log.Error().Err(err).Msg("client.GetAllProjects(ctx, collection)")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
