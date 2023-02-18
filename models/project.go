package models

type Project struct {
	Id           string   `json:"id" firestore:"id"`
	Name         string   `json:"name" firestore:"name"`
	Environments []string `json:"environments,omitempty" firestore:"environments,omitempty"`
}
