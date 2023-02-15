package gcp

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
)

type Client struct {
	Firestore *firestore.Client
}

func New(projectId string) (*Client, error) {
	ctx := context.Background()

	c, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("error creating firestore client : %v", err)
	}

	return &Client{c}, err
}
