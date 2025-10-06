# Atlas Go SDK Examples for the MongoDB Atlas Architecture Center

This repository contains runnable examples for the
[Atlas Go SDK](https://www.mongodb.com/docs/atlas/sdk/)
that align with best practices from the MongoDB
[Atlas Architecture Center](https://www.mongodb.com/docs/atlas/architecture/current/).

Use these examples as starting points for your own Atlas integration.

## Features

Currently, the repository includes examples that demonstrate the following:

- Authenticate with service accounts
- Return cluster and database metrics
- Download logs for a specific host
- Pull and parse line-item-level billing data
- Return all linked organizations from a specific billing organization
- Get historical invoices for an organization
- Programmatically archive Atlas cluster data
- Proactively or reactively scale clusters based on configuration

As the Architecture Center documentation evolves, this repository will be updated with new examples 
and improvements to existing code. 

## Project Structure

```text
.
├── examples             # Runnable examples by category
│   ├── billing/
│   ├── monitoring/
│   └── performance/
├── configs              # Atlas configuration templates & environment-specific configs
│   └── config.example.json
├── internal             # Shared utilities and helpers
│   ├── archive/
│   ├── auth/
│   ├── billing/
│   ├── clusterutils/
│   ├── config/
│   ├── data/
│   ├── errors/
│   ├── fileutils/
│   ├── logs/
│   ├── metrics/
│   └── scale/
├── go.mod
├── go.sum
├── CHANGELOG.md         # List of major changes to the project 
├── .gitignore           # Ignores .env file and log output
└── .env.example         # Example environment variables
```

## Prerequisites

- Go 1.24 or later
- A MongoDB Atlas organization, project, and at least one cluster
- Service account credentials with appropriate permissions and IP access. See
  [Service Account Overview](https://www.mongodb.com/docs/atlas/api/service-accounts-overview/)

## Environment Variables

Only a small set of environment variables are required. Programmatic scaling and DR settings are provided via the JSON config file — not separate env vars.

Create a `.env.<environment>` file (e.g. `.env.development`):

```dotenv
# Required service account credentials
MONGODB_ATLAS_SERVICE_ACCOUNT_ID=<your_service_account_id>
MONGODB_ATLAS_SERVICE_ACCOUNT_SECRET=<your_service_account_secret>

# Optional: override default config path (defaults to configs/config.json if unset)
CONFIG_PATH=configs/config.development.json

# Optional: base directory for downloaded artifacts (logs, archives, invoices)
ATLAS_DOWNLOADS_DIR=tmp/atlas_downloads
```

> NOTE: For production, store secrets in a secrets manager (e.g. HashiCorp Vault, AWS Secrets Manager) instead of plain environment variables. See [Secrets management](https://www.mongodb.com/docs/atlas/architecture/current/auth/#secrets-management).

## Configuration File

Create `configs/config.<environment>.json` (e.g. `configs/config.development.json`). If `CONFIG_PATH` is unset, the loader falls back to `configs/config.json`.

Minimal example:
```json
{
  "ATLAS_ORG_ID": "<your-org-id>",
  "ATLAS_PROJECT_ID": "<your-project-id>",
  "ATLAS_CLUSTER_NAME": "<a-cluster-name>",
  "ATLAS_PROCESS_ID": "<cluster-hostname:port>",
  "programmatic_scaling": {
    "target_tier": "M50",
    "pre_scale_event": false,
    "cpu_threshold": 75.0,
    "cpu_period_minutes": 60,
    "dry_run": true
  }
}
```

Field notes:
- `ATLAS_PROCESS_ID` is used for examples that operate directly on a single host (logs/metrics). Format: `hostname:port`.
- `programmatic_scaling` (optional) controls proactive (pre_scale_event) and reactive (cpu_threshold over cpu_period_minutes) scaling.
- `dry_run=true` ensures scaling logic logs intent without applying changes.
- Omit `programmatic_scaling` entirely to skip scaling analysis.
- Omit `disaster_recovery` if not exercising DR examples.

Defaults applied when absent:
- `programmatic_scaling.target_tier` → `M50`
- `programmatic_scaling.cpu_threshold` → `75.0`
- `programmatic_scaling.cpu_period_minutes` → `60`
- `programmatic_scaling.dry_run` → `true`

## Running Examples

Each example is an independent entrypoint. Ensure your `.env.<env>` and matching config file are in place, then:

```bash
# Example: run with development environment
cp .env.example .env.development   # (or create manually)
# edit .env.development and config file with real values

# Billing - historical invoices
go run examples/billing/historical/main.go

# Billing - line items
go run examples/billing/line_items/main.go

# Billing - linked organizations
go run examples/billing/linked_orgs/main.go

# Logs - fetch host logs
go run examples/monitoring/logs/main.go

# Metrics - disk measurements
go run examples/monitoring/metrics_disk/main.go

# Metrics - process CPU metrics
go run examples/monitoring/metrics_process/main.go

# Performance - archive cluster data
go run examples/performance/archiving/main.go

# Performance - programmatic scaling (dry run by default)
go run examples/performance/scaling/main.go
```

### Programmatic Scaling Behavior

The scaling example evaluates each cluster:
1. Skips non-IDLE clusters.
2. Applies `pre_scale_event` first (immediate scale intent).
3. For dedicated tiers: collects per-process CPU, prioritizes primary; falls back to aggregated average across processes.
4. For shared tiers (M0/M2/M5): skips reactive CPU (metrics limited); only pre-scale can trigger.
5. When `dry_run=false`, executes a tier change to `target_tier`.

## Changelog

For a list of major changes to this project, see [CHANGELOG](CHANGELOG.md).

## Reporting Issues

Use the "Rate this page" widget on the
[Atlas Architecture Center](https://www.mongodb.com/docs/atlas/architecture/current/)
docs to leave feedback or file issues.

## License

This project is licensed under Apache 2.0. See [LICENSE](LICENSE.md).
