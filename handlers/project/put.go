package project

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

func Put(client *firestore.Client, config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data models.Project
		id := c.Param("id")

		if err := c.BindJSON(&data); err != nil {
			log.Error().Err(err).Msg("error binding request body to struct")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("error binding request body to struct : %v", err),
			})
			return
		}

		if _, err := client.Collection(config.Collection).Doc(id).Set(c, data); err != nil {
			log.Error().Err(err).Msg("error updating document in collection")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("error updating document in collection : %v", err),
			})
			return
		}

		for _, env := range data.Environments {
			stackName, err := iac.Up(c, data, env, config)
			if err != nil {
				log.Error().Err(err).Msg("error updating environment")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("error updating environment : %v", err),
				})
				return
			}

			if _, err := client.Collection(config.Collection).Doc(id).Update(c, []firestore.Update{
				{
					Path:  "stackName",
					Value: stackName,
				}}); err != nil {
				log.Error().Err(err).Msg("error updating document with stack name")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("error updating document with stack name : %v", err),
				})
				return
			}
		}

		c.Status(http.StatusAccepted)
	}
}
