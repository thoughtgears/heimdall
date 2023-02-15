package models

type Project struct {
	Name         string   `json:"name" firestore:"name"`
	Environments []string `json:"environments" firestore:"environments"`
}
