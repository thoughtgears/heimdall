package models

type Project struct {
	Id           string   `json:"id" firestore:"id"`
	Name         string   `json:"name" firestore:"name"`
	Repository   string   `json:"repository" firestore:"repository"`
	Services     []string `json:"services" firestore:"services"`
	Environments []string `json:"environments,omitempty" firestore:"environments,omitempty"`
}
