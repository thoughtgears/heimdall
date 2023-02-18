package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"os"
)

func main() {
	projectId := os.Getenv("GOOGLE_PROJECT")

	pulumi.Run(func(ctx *pulumi.Context) error {
		project, err := organizations.LookupProject(ctx, &organizations.LookupProjectArgs{ProjectId: &projectId})
		if err != nil {
			return err
		}
		return nil
	})
}
