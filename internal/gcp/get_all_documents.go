package gcp

import (
	"context"
	"fmt"
	"github.com/thoughtgears/heimdall/models"
	"google.golang.org/api/iterator"
)

func (c *Client) GetAllProjects(ctx context.Context, collection string) ([]models.Project, error) {
	var resp []models.Project
	iter := c.Firestore.Collection(collection).Documents(ctx)

	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading document from collection : %v", err)
		}
		var p models.Project
		if err := doc.DataTo(&p); err != nil {
			return nil, fmt.Errorf("error binding document to struct : %v", err)
		}
		resp = append(resp, p)
	}

	return resp, nil
}
