package models

import "time"

type ServiceCatalogue struct {
	Id               string    `json:"id" firestore:"id"`
	Services         []string  `json:"services" firestore:"services"`
	CatalogueVersion float32   `json:"catalogueVersion" firestore:"catalogueVersion"`
	LastUpdate       time.Time `json:"lastUpdate,omitempty" firestore:"lastUpdate,omitempty"`
}
