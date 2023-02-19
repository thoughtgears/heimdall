package main

import (
	"fmt"
	"heimdall/core"

	"github.com/kelseyhightower/envconfig"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/firestore"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
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
		if err := core.WorkloadIdentityFederation(ctx, env.ProjectId, env.WorkloadPool, serviceName, githubIssuerUrl, env.Repository); err != nil {
			return err
		}

		repository, err := core.Repository(ctx, serviceName, env.ProjectId, env.Region)
		if err != nil {
			return nil
		}

		cloudRunSA, err := serviceaccount.NewAccount(ctx, fmt.Sprintf("cloud-run-%s", serviceName), &serviceaccount.AccountArgs{
			AccountId:   pulumi.Sprintf("cloud-run-%s", serviceName),
			DisplayName: pulumi.Sprintf("Cloud run service account for %s", serviceName),
			Project:     pulumi.String(env.ProjectId),
		})
		if err != nil {
			return err
		}

		if err := core.RepositoryPermissions(ctx, cloudRunSA, repository); err != nil {
			return err
		}

		//_, err := cloudrun.NewService(ctx, serviceName, &cloudrun.ServiceArgs{
		//	Location: pulumi.String(env.Region),
		//	Template: &cloudrun.ServiceTemplateArgs{
		//		Spec: &cloudrun.ServiceTemplateSpecArgs{
		//			Containers: cloudrun.ServiceTemplateSpecContainerArray{
		//				&cloudrun.ServiceTemplateSpecContainerArgs{
		//					Image: pulumi.Sprintf("%s-docker.pkg.dev/%s/%s", env.Region, repository.Name, serviceName),
		//				},
		//			},
		//		},
		//	},
		//	Traffics: cloudrun.ServiceTrafficArray{
		//		&cloudrun.ServiceTrafficArgs{
		//			LatestRevision: pulumi.Bool(true),
		//			Percent:        pulumi.Int(100),
		//		},
		//	},
		//})
		//if err != nil {
		//	return err
		//}

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

		return nil
	})
}
