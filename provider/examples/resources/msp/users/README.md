# MSP users Example

This example shows you how to add users to an MSP managed tenant. 

## Pre-requisites

You need access to an MSP Portal, and API token for the MSP portal.

## Usage
- Modify `providers.tf` accordingly.
- Paste CDO API token for an MSP portal into `api_token.txt`
    - see https://docs.defenseorchestrator.com/#!c-api-tokens.html for how to generate this.
- Specify the name of a tenant managed by the MSP Portal. You can get the tenant name by going to Settings in the MSP portal.
- To see the generated API token for the created user, run `terraform show -json | jq -r ".values.outputs.api_token.value"`