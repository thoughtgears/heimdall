package project

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/thoughtgears/heimdall/models"
	"google.golang.org/api/iterator"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/thoughtgears/heimdall/internal/config"
)

func GetAll(client *firestore.Client, config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp []models.Project
		iter := client.Collection(config.Collection).Documents(c)

		defer iter.Stop()

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Error().Err(err).Msg("error reading document from collection")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("error reading document from collection : %v", err),
				})
				return
			}
			var p models.Project
			if err := doc.DataTo(&p); err != nil {
				log.Error().Err(err).Msg("error binding document to struct")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": fmt.Sprintf("error binding document to struct : %v", err),
				})
				return
			}
			resp = append(resp, p)
		}

		c.JSON(http.StatusOK, resp)
	}
}
