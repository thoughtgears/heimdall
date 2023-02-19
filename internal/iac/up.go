package iac

import (
	"context"
	"fmt"
	"os"

	"github.com/thoughtgears/heimdall/internal/config"
	"github.com/thoughtgears/heimdall/models"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Up(ctx context.Context, data models.Project, environment string, config *config.Config) (string, error) {
	deploy := func(ctx *pulumi.Context) error {
		// Ensure you have a proper suffix ID at the end of creating a project
		projectSuffix, err := random.NewRandomInteger(ctx, "project_suffix", &random.RandomIntegerArgs{
			Min: pulumi.Int(10000),
			Max: pulumi.Int(99999),
		})
		if err != nil {
			return err
		}

		// Create the new landing-zone project to have as a baseline when using pulumi to be able to create other
		// projects in your GCP organization.
		project, err := organizations.NewProject(ctx, data.Name, &organizations.ProjectArgs{
			Name:              pulumi.String(data.Name),
			BillingAccount:    pulumi.String(data.BillingAccount),
			AutoCreateNetwork: pulumi.Bool(false),
			OrgId:             pulumi.String(data.OrganizationId),
			ProjectId:         pulumi.Sprintf("%s-%v", data.Name, projectSuffix.Result),
		})
		if err != nil {
			return err
		}

		ctx.Export("project id", project.ProjectId)

		return nil
	}

	stackName := auto.FullyQualifiedStackName(data.PulumiOwner, data.Name, environment)
	s, err := auto.UpsertStackInlineSource(ctx, stackName, data.Name, deploy)
	w := s.Workspace()

	// for inline source programs, we must manage plugins ourselves
	err = w.InstallPlugin(ctx, "gcp", "v6.50.0")
	if err != nil {
		return "", fmt.Errorf("failed to install gcp plugin : %v", err)
	}

	// set stack configuration specifying the AWS region to deploy
	if err := s.SetAllConfig(ctx, auto.ConfigMap{
		"gcp:region":  auto.ConfigValue{Value: config.Region},
		"gcp:project": auto.ConfigValue{Value: config.Project},
	}); err != nil {
		return "", fmt.Errorf("failed setting config : %v", err)
	}

	_, err = s.Refresh(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to refresh stack : %v", err)
	}

	// wire up our update to stream progress to stdout
	stdoutStreamer := optup.ProgressStreams(os.Stdout)

	if _, err := s.Up(ctx, stdoutStreamer); err != nil {
		return "", fmt.Errorf("failed to update stack : %v", err)
	}

	return stackName, nil
}
