package project

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/thoughtgears/heimdall/internal/config"
)

func Delete(client *firestore.Client, config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if _, err := client.Collection(config.Collection).Doc(id).Delete(c); err != nil {
			log.Error().Err(err).Msg("error removing document from collection")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("error removing document from collection : %v", err),
			})
			return
		}

		c.Status(http.StatusAccepted)
	}
}
