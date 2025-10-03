package clusterutils

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250219001/admin"

	"atlas-sdk-go/internal/errors"
)

// GetClusterSRVConnectionString returns the standard SRV connection string for a cluster.
func GetClusterSRVConnectionString(ctx context.Context, client *admin.APIClient, projectID, clusterName string) (string, error) {
	if client == nil {
		return "", fmt.Errorf("nil atlas api client")
	}
	cluster, _, err := client.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		return "", errors.FormatError("get cluster", projectID, err)
	}
	if cluster == nil || cluster.ConnectionStrings == nil || cluster.ConnectionStrings.StandardSrv == nil {
		return "", fmt.Errorf("no standard SRV connection string found for cluster %s", clusterName)
	}
	return *cluster.ConnectionStrings.StandardSrv, nil
}
