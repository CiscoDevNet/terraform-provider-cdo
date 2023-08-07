### CDO Terraform Provider

This repo provides a terraform provider to provision resources on the Cisco CDO platform.

#### Structure

We make use of Go workspaces - this is to split the repsonsibility of the provider and the Go CDO client. 
Eventually the client will be moved to its own repo, but in the interest of hitting the ground running they both live in here.

```
.
├── client        # Golang CDO client - currently using the module name github.com/CiscoDevnet/terraform-provider-cdo/go-client
├── provider      # Terraform provider
└── README.md
```

#### Requirements

* Go (1.20)
  - macos install: `brew install go`
* tfenv 
  - macos install: `brew install tfenv`

#### Acceptance Tests

**Acceptance tests will create real resources!**

Ensure you have met the requirements.

Ensure you have an active terraform:

```bash
tfenv use 1.3.1
```

Then in the provider dir you can run the acceptances tests:

```bash
cd ./provider
ACC_TEST_CISCO_CDO_API_TOKEN=<CDO_API_TOKEN> make testacc
```

### Running Examples
Examples are provided so that you can do the usual `plan`, `apply`, `destroy` etc in the folders under `provider/examples` directory.
#### Setup
1. Build provider locally.
   1. Figure out where your executable will be built.
      ```bash
      go env GOBIN
      ```
      **Note:** If empty output, it is default to `~/go/bin`
   2. Build executable.
      ```bash
      cd provider
      go install .
      ```
   3. Verify installation
      The installation path that you figured out above should contains an executable called **terraform-provider-cdo**.
2. Tell terraform to use your local build
   1. Modify local terraform configuration
      ```bash
      vim ~/.terraformrc
      ```
   2. Copy the following into it:
      ```terraform
      provider_installation {
   
        dev_overrides {
            "hashicorp.com/CiscoDevnet/cdo" = "<DIRECTORY_OF_YOUR_EXECUTABLE>"
        }
  
        # For all other providers, install them directly from their origin provider
        # registries as normal. If you omit this, Terraform will _only_ use
        # the dev_overrides block, and so no other providers will be available.
        direct {}
      }
      ```
#### Running
1. Navigate to a folder under `provider/examples`, e.g. `provider/examples/resources/asa`.
2. Run
   ```bash
   terraform plan
   ```
3. Output (sample)
   ```
    |
    │ Warning: Provider development overrides are in effect
    │
    │ The following provider development overrides are set in the CLI configuration:
    │  - hashicorp.com/CiscoDevnet/cdo in /Users/weilluo/go/bin
    │
    │ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
    |
    ... rest
   ```
#### Troubleshooting
- Error: Inconsistent dependency lock file
  ```
  provider hashicorp/CiscoDevnet/cdo: required by this configuration but no version is selected
  ```
  - This means you have not setup the dev override properly, make sure your `~/.terraformrc` has the right override for the provider in question.