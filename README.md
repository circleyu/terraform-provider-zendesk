# terraform-provider-zendesk

[![Actions Status](https://github.com/nukosuke/terraform-provider-zendesk/workflows/CI/badge.svg)](https://github.com/nukosuke/terraform-provider-zendesk/actions)
[![Build status](https://ci.appveyor.com/api/projects/status/ti5il35v6a6ankcq/branch/master?svg=true)](https://ci.appveyor.com/project/nukosuke/terraform-provider-zendesk/branch/master)
[![Coverage Status](https://coveralls.io/repos/github/nukosuke/terraform-provider-zendesk/badge.svg?branch=master)](https://coveralls.io/github/nukosuke/terraform-provider-zendesk?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/nukosuke/terraform-provider-zendesk)](https://goreportcard.com/report/github.com/nukosuke/terraform-provider-zendesk)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fnukosuke%2Fterraform-provider-zendesk.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fnukosuke%2Fterraform-provider-zendesk?ref=badge_shield)

# circleyu Additions
14.12.25 - Added support for new Zendesk resources and data sources:
  - Resources: Custom Roles, Custom Statuses, Users, Tickets, Queues, Group Memberships, Organization Memberships
  - Data Sources: Tags, Locales, Satisfaction Ratings
  - Fix linter warnings and update deprecated APIs (rand.Seed, io/ioutil, SchemaValidateFunc)
09.05.25 - Fix the error in the Target update.

# clearnote01 Additions

20.03.24 - Added support for Zendesk Views, Macros & User Fields
03.04.23 - Added support for Zendesk Trigger Categories

# Debugging

cd examples_scratchpad
export TF_LOG=trace/debug/info
make changes in scratchpad/resource.tf file in this folder
make plan/apply

Terraform provider for Zendesk

- [Document](https://registry.terraform.io/providers/circleyu/zendesk/latest/docs)

## Requirements

- Terraform >= v0.12.7
- Go >= v1.18 (only for build)

## Usage

terraform-provider-zendesk is available on [Terraform Registry](https://registry.terraform.io). You don't need to download artifacts manually.  
Instead configure provider as follow.

```hcl
terraform {
  required_providers {
    zendesk = {
      source  = "circleyu/zendesk"
      version = ">= 0.0.2"
    }
  }
}
```

and run `terraform init` in your Terraform resource directory.

## Development

### Build from source

```sh
$ git clone git@github.com:clearnote01/terraform-provider-zendesk.git
$ cd terraform-provider-zendesk
$ go mod download
$ go build
```

## License

MIT License

See the file [LICENSE](./LICENSE) for details.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fnukosuke%2Fterraform-provider-zendesk.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fnukosuke%2Fterraform-provider-zendesk?ref=badge_large)

## Related project

- [nukosuke/go-zendesk](https://github.com/nukosuke/go-zendesk)
