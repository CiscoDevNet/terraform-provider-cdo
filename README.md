### CDO Terraform Provider

This repo provides a terraform provider to provision resources on the Cisco CDO platform.

#### Structure

We make use of Go workspaces - this is to split the repsonsibility of the provider and the Go CDO client. 
Eventually the client will be moved to its own repo, but in the interest of hitting the ground running they both live in here.

```
.
├── client        # Golang CDO client - currently using the module name github.com/cisco-lockhart/go-client
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