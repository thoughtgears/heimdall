package service_cataloge

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/thoughtgears/heimdall/models"
)

func Get(client *firestore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data models.ServiceCatalogue

		doc, err := client.Collection("pulumiData").Doc("serviceCatalogue").Get(c)
		if err != nil {
			log.Error().Err(err).Msg("error retrieving document from firestore")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("error retrieving document from firestore : %v", err),
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

		c.JSON(http.StatusOK, data)
	}
}
