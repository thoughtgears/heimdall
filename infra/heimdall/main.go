package main

import (
	"github.com/kelseyhightower/envconfig"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/iam"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/firestore"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const serviceName = "heimdall"
const githubIssuerUrl = "https://token.actions.githubusercontent.com"

type Environment struct {
	ProjectId    string `envconfig:"GOOGLE_PROJECT" required:"true"`
	Region       string `envconfig:"REGION" default:"europe-west2"`
	WorkloadPool string `envconfig:"WORKLOAD_POOL" default:"pulimi-landingzone"`
	Repository   string `envconfig:"GIT_REPO" required:"true"`
}

var env Environment

func init() {
	envconfig.MustProcess("", &env)
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		pool, err := iam.GetWorkloadIdentityPool(ctx, env.WorkloadPool, pulumi.ID(env.WorkloadPool), nil)
		if err != nil {
			return err
		}

		p, err := iam.NewWorkloadIdentityPoolProvider(ctx, serviceName, &iam.WorkloadIdentityPoolProviderArgs{
			WorkloadIdentityPoolId:         pool.WorkloadIdentityPoolId,
			WorkloadIdentityPoolProviderId: pulumi.String(serviceName),
			AttributeMapping: pulumi.StringMap{
				"google.subject":       pulumi.String("assertion.sub"),
				"attribute.actor":      pulumi.String("assertion.actor"),
				"attribute.repository": pulumi.String("assertion.repository"),
			},
			Project: pulumi.String(env.ProjectId),
			Oidc: &iam.WorkloadIdentityPoolProviderOidcArgs{
				IssuerUri: pulumi.String(githubIssuerUrl),
			},
		})
		if err != nil {
			return err
		}

		// Set the Workload provider to interact with the provider through our pulumi SA
		if _, err := serviceaccount.NewIAMMember(ctx, "heimdall-workload-identity", &serviceaccount.IAMMemberArgs{
			Member:           pulumi.Sprintf("principalSet://iam.googleapis.com/%s/attribute.repository/%s", pool.Name, env.Repository),
			Role:             pulumi.String("roles/iam.workloadIdentityUser"),
			ServiceAccountId: pulumi.Sprintf("projects/%s/serviceAccounts/pulumi@%s.iam.gserviceaccount.com", env.ProjectId, env.ProjectId),
		}, pulumi.DependsOn([]pulumi.Resource{p})); err != nil {
			return err
		}

		database, err := firestore.NewDatabase(ctx, "heimdall", &firestore.DatabaseArgs{
			AppEngineIntegrationMode: pulumi.String("DISABLED"),
			ConcurrencyMode:          pulumi.String("PESSIMISTIC"),
			LocationId:               pulumi.String(env.Region),
			Name:                     pulumi.String(`(default)`),
			Project:                  pulumi.String(env.ProjectId),
			Type:                     pulumi.String("FIRESTORE_NATIVE"),
		})
		if err != nil {
			return err
		}

		ctx.Export("firestore database", database.ID())
		ctx.Export("workload-identity-pool-id:", p.ID())

		return nil
	})
}
