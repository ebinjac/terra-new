# Terraform CLI Tool

This CLI tool automates the process of downloading Terraform configurations, generating graphs, and uploading plan data.

## Prerequisites

- Go 1.21 or later
- Terraform installed and available in PATH
- `iriusrisk.yaml` file in the project root directory

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd terraform-cli
```

2. Build the project:
```bash
go build -o tfcli ./cmd/tfcli
```

## Usage

1. Create an input JSON file with the required configuration:

```json
{
  "config_url": "https://example.com/api/v2/config/download",
  "plan_url": "https://example.com/api/v2/plans/plan-id/json-output",
  "token": "your-bearer-token",
  "product-id": "your-product-id",
  "name": "project-name"
}
```

2. Run the tool:
```bash
./tfcli input.json
```

The tool will:
1. Download and extract the Terraform configuration
2. Run `terraform init`
3. Generate a Terraform graph
4. Download the plan.json file
5. Upload all required files to the specified endpoint

## Error Handling

The tool includes comprehensive error handling and will exit with a non-zero status code if any step fails. Error messages will be printed to stderr with details about what went wrong. 