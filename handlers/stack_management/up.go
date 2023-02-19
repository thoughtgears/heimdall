package stack_management

import (
	"fmt"
	"net/http"

	"github.com/thoughtgears/heimdall/internal/iac"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/thoughtgears/heimdall/internal/config"
	"github.com/thoughtgears/heimdall/models"
)

func Up(client *firestore.Client, config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data models.Project
		id := c.Param("id")

		doc, err := client.Collection(config.Collection).Doc(id).Get(c)
		if err != nil {
			log.Error().Err(err).Msg("error reading document from collection")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("error reading document from collection : %v", err),
			})
			return
		}

		if err := doc.DataTo(&data); err != nil {
			log.Error().Err(err).Msg("error binding document to struct")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("error binding document to struct : %v", err),
			})
			return
		}

		for _, env := range data.Environments {
			if _, err := iac.Up(c, data, env, config); err != nil {
				log.Error().Err(err).Msg("error updating environment")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("error updating environment : %v", err),
				})
				return
			}
		}

		c.Status(http.StatusAccepted)
	}
}
