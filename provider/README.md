# Terraform Provider Scaffolding (Terraform Plugin Framework)

_This template repository is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). The template repository built on the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) can be found at [terraform-provider-scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding). See [Which SDK Should I Use?](https://www.terraform.io/docs/plugin/which-sdk.html) in the Terraform documentation for additional information._

This repository is a _template_ for a [Terraform](https://www.terraform.io) provider. It is intended as a starting point for creating Terraform providers, containing:

- A resource and a data source (`internal/provider/`),
- Examples (`examples/`) and generated documentation (`docs/`),
- Miscellaneous meta files.

These files contain boilerplate code that you will need to edit to create your own Terraform provider. Tutorials for creating Terraform providers can be found on the [HashiCorp Learn](https://learn.hashicorp.com/collections/terraform/providers-plugin-framework) platform. _Terraform Plugin Framework specific guides are titled accordingly._

Please see the [GitHub template repository documentation](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/creating-a-repository-from-a-template) for how to create a new repository from this template on GitHub.

Once you've written your provider, you'll want to [publish it on the Terraform Registry](https://www.terraform.io/docs/registry/providers/publishing.html) so that others can use it.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Fill this in for each provider

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

### Testing

Terraform providers generally use acceptance tests to verify behaviours, so it will creates, destroy and update real infrastructure. This means that the tests will require to interact with CDO, so it needs some setups.

The acceptance tests are run in the CDO CI environment. You need to be part of the CDO team to access this tenant.

#### API Token

In order to interact with CDO, an API-only user with super-admin privileges exists in the tenant, and a token generated for that user is used to run the acceptance tests. If you are a CDO engineer with access to this tenant, please log in as a super-admin user on this tenant and create an API-only user for your local testing. If you are not a CDO engineer, you will have to rely on the pull requests you create to be tested on Github.

When developing locally, you can overwrite this token using environment variable, see files in `provider/internal/acctest` for details.

## CI

### Devices

Since terraform requires interacting with real infrastructure, tests under the `provider` folder are mostly acceptance tests, and they communicate with devices in the CI environment. At the moment we need to create and onboard some devices in CI for our tests to work. (It will be more self-contained in the future, and this section will no longer be needed, hopefully).