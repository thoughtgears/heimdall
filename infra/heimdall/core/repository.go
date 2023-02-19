package core

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Repository(ctx *pulumi.Context, serviceName, projectId, region string) (*artifactregistry.Repository, error) {
	api, err := projects.NewService(ctx, "artifactregistry.googleapis.com", &projects.ServiceArgs{
		DisableDependentServices: pulumi.Bool(true),
		Project:                  pulumi.String(projectId),
		Service:                  pulumi.String("artifactregistry.googleapis.com"),
	})
	if err != nil {
		return nil, err
	}

	r, err := artifactregistry.NewRepository(ctx, serviceName, &artifactregistry.RepositoryArgs{
		Description:  pulumi.Sprintf("repository for %s", serviceName),
		Format:       pulumi.String("DOCKER"),
		Location:     pulumi.String(region),
		RepositoryId: pulumi.String(serviceName),
	}, pulumi.DependsOn([]pulumi.Resource{api}))
	if err != nil {
		return nil, err
	}

	ctx.Export("artifact repository", r.RepositoryId)

	return r, nil
}

func RepositoryPermissions(ctx *pulumi.Context, serviceAccount *serviceaccount.Account, repository *artifactregistry.Repository) error {
	if _, err := artifactregistry.NewRepositoryIamMember(ctx, fmt.Sprintf("iam-member-%v", serviceAccount.Name.ApplyT(
		func(arg interface{}) (string, error) {
			return arg.(string), nil
		})), &artifactregistry.RepositoryIamMemberArgs{
		Project:    repository.Project,
		Location:   repository.Location,
		Repository: repository.Name,
		Role:       pulumi.String("roles/artifactregistry.reader"),
		Member:     pulumi.Sprintf("serviceAccount:%s", serviceAccount.Email),
	}); err != nil {
		return err
	}

	return nil
}
