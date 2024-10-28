# MSP tenants Example

## Pre-requisites

You need access to an MSP Portal, and API token for the MSP portal.

## Usage
- Modify `terraform.tfvars` and `providers.tf` accordingly.
- Paste CDO API token for an MSP portal into `api_token.txt`
    - see https://docs.defenseorchestrator.com/#!c-api-tokens.html for how to generate this.
- Specify the name of a tenant managed by the MSP Portal. You can get the tenant name by going to Settings in the MSP portal.