package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/firestore"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"os"
)

func main() {
	projectId := os.Getenv("GOOGLE_PROJECT")
	region := os.Getenv("REGION")

	pulumi.Run(func(ctx *pulumi.Context) error {
		project, err := organizations.LookupProject(ctx, &organizations.LookupProjectArgs{ProjectId: &projectId})
		if err != nil {
			return err
		}

		database, err := firestore.NewDatabase(ctx, "heimdall", &firestore.DatabaseArgs{
			AppEngineIntegrationMode: pulumi.String("ENABLED"),
			ConcurrencyMode:          pulumi.String("PESSIMISTIC"),
			LocationId:               pulumi.String(region),
			Name:                     pulumi.String(`(default)`),
			Project:                  pulumi.String(*project.ProjectId),
			Type:                     pulumi.String("FIRESTORE_NATIVE"),
		})
		if err != nil {
			return err
		}

		ctx.Export("firestore database", database.Name)

		return nil
	})
}
