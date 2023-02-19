package core

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/iam"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func WorkloadIdentityFederation(ctx *pulumi.Context, projectId, workloadPool, serviceName, issuerUrl, repository string) error {
	pool, err := iam.GetWorkloadIdentityPool(ctx, workloadPool, pulumi.ID(workloadPool), nil)
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
		Project: pulumi.String(projectId),
		Oidc: &iam.WorkloadIdentityPoolProviderOidcArgs{
			IssuerUri: pulumi.String(issuerUrl),
		},
	})
	if err != nil {
		return err
	}

	// Set the Workload provider to interact with the provider through our pulumi SA
	if _, err := serviceaccount.NewIAMMember(ctx, "heimdall-workload-identity", &serviceaccount.IAMMemberArgs{
		Member:           pulumi.Sprintf("principalSet://iam.googleapis.com/%s/attribute.repository/%s", pool.Name, repository),
		Role:             pulumi.String("roles/iam.workloadIdentityUser"),
		ServiceAccountId: pulumi.Sprintf("projects/%s/serviceAccounts/pulumi@%s.iam.gserviceaccount.com", projectId, projectId),
	}, pulumi.DependsOn([]pulumi.Resource{p})); err != nil {
		return err
	}

	ctx.Export("workload-identity-pool-id:", p.ID())

	return nil
}
