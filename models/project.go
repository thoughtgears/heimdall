package models

import (
	"crypto/md5"
	"encoding/hex"
)

type Project struct {
	Name            string   `json:"name" firestore:"name"`
	BillingAccount  string   `json:"billingAccount" firestore:"billingAccount"`
	OrganizationId  string   `json:"organizationId" firestore:"organizationId"`
	Repository      string   `json:"repository" firestore:"repository"`
	RepositoryOwner string   `json:"repositoryOwner" firestore:"repositoryOwner"`
	PulumiOwner     string   `json:"pulumiOwner" firestore:"pulumiOwner"`
	StackName       string   `json:"stackName" firestore:"stackName"`
	Services        []string `json:"services" firestore:"services"`
	Environments    []string `json:"environments,omitempty" firestore:"environments,omitempty"`
}

// GetProjectId takes the name and makes a md5 hash to return to have as ID
func GetProjectId(id string) string {
	hash := md5.Sum([]byte(id))
	return hex.EncodeToString(hash[:])
}
