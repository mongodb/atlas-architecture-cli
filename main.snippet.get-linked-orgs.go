// See entire project at https://github.com/mongodb/atlas-architecture-go-sdk
package main

import (
	"context"
	"fmt"
	"log"

	"atlas-sdk-go/internal/auth"
	"atlas-sdk-go/internal/billing"
	"atlas-sdk-go/internal/config"

	"github.com/joho/godotenv"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

func main() {
	_ = godotenv.Load()

	secrets, cfg, err := config.LoadAll("configs/config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	sdk, err := auth.NewClient(cfg, secrets)
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	ctx := context.Background()
	params := &admin.ListInvoicesApiParams{
		OrgId: cfg.OrgID,
	}

	fmt.Printf("Fetching cross-org billing info for organization: %s\n", params.OrgId)
	results, err := billing.GetCrossOrgBilling(ctx, sdk.InvoicesApi, params)
	if err != nil {
		log.Fatalf("Failed to retrieve invoices: %v", err)
	}
	if len(results) == 0 {
		fmt.Println("No invoices found for the billing organization")
		return
	}

	linkedOrgs, err := billing.GetLinkedOrgs(ctx, sdk.InvoicesApi, params)
	if err != nil {
		log.Fatalf("Failed to retrieve linked organizations: %v", err)
	}
	if len(linkedOrgs) == 0 {
		fmt.Println("No linked organizations found for the billing org")
		return
	}
	fmt.Println("Linked organizations:")
	for i, org := range linkedOrgs {
		fmt.Printf("  %d. %v\n", i+1, org)
	}
}
