package testimpl

import (
	"context"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	armPostgres "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresqlflexibleservers"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/launchbynttdata/lcaf-component-terratest/types"
	"github.com/stretchr/testify/assert"
)

func TestPostgresqlServer(t *testing.T, ctx types.TestContext) {
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")
	if len(subscriptionId) == 0 {
		t.Fatal("ARM_SUBSCRIPTION_ID environment variable is not set")
	}

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		t.Fatalf("Unable to get credentials: %e\n", err)
	}

	options := arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloud.AzurePublic,
		},
	}

	armPostgresConfigClient, err := armPostgres.NewConfigurationsClient(subscriptionId, credential, &options)
	if err != nil {
		t.Fatalf("Error getting Postgres client: %v", err)
	}

	t.Run("doesPostgresqlServerConfigurationExist", func(t *testing.T) {
		resourceGroupName := terraform.Output(t, ctx.TerratestTerraformOptions(), "resource_group_name")
		postgresName := terraform.Output(t, ctx.TerratestTerraformOptions(), "postgres_server_name")
		postgresConfig := terraform.OutputMap(t, ctx.TerratestTerraformOptions(), "postgres_server_config")

		for key := range postgresConfig {
			postgresqlConfig, err := armPostgresConfigClient.Get(context.Background(), resourceGroupName, postgresName, key, nil)
			if err != nil {
				t.Fatalf("Error getting Postgresql server: %v", err)
			}

			assert.Equal(t, key, *postgresqlConfig.Name)
		}
	})
}
