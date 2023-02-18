package models

import "time"

type ServiceCatalogue struct {
	Services         []string  `json:"services" firestore:"services"`
	CatalogueVersion int       `json:"catalogueVersion" firestore:"catalogueVersion"`
	LastUpdate       time.Time `json:"lastUpdate,omitempty" firestore:"lastUpdate,omitempty"`
}
