package service_cataloge

import (
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/thoughtgears/heimdall/models"
)

func Post(client *firestore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data models.ServiceCatalogue

		if err := c.BindJSON(&data); err != nil {
			log.Error().Err(err).Msg("error binding request body to struct")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("error binding request body to struct : %v", err),
			})
			return
		}

		data.LastUpdate = time.Now()

		if _, err := client.Collection("pulumiData").Doc("serviceCatalogue").Set(c, data); err != nil {
			log.Error().Err(err).Msg("error writing document to firestore")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("error writing document to firestore : %v", err),
			})
			return
		}

		c.Status(http.StatusAccepted)
	}
}
