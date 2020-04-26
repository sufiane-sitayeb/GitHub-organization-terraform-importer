# GitHub Organization Terraform Importer

![License](https://img.shields.io/github/license/AlessioCasco/GitHub-organization-terraform-importer)

[![Go Report Card](https://goreportcard.com/badge/github.com/AlessioCasco/GitHub-organization-terraform-importer)](https://goreportcard.com/report/github.com/AlessioCasco/GitHub-organization-terraform-importer)

This tool fetches some of your resources from GitHub organization (repos, teams and users) and generates the correspondent terraform files and modules.
It also allows you to import the resources into a terraform state so you can start to manage them with terraform itself.

## Disclaimer

1) Why GO: This project was a way for me to start playing and learn it.

2) Your code is quite bad...: See point 1 and feel free to open a PR with improvement or an issue with suggestions, I'll try to improve it from time to time.

## Prerequisites

* [GitHub Personal access tokens](https://GitHub.com/settings/tokens) with the following selected:
  * repo section:
    * Select everything
  * user section:
    * `read:user`
  * Admin:org section:
    * `read:org`

* A GitHub organization
* Go (tested with version 1.14)
* Terraform >= v0.12.6 (modules use `for_each`)
* If you set a backend file for terraform, you need the credentials file that allows writing on the backend. (S3 for example)

## Parameters

`-import`: Imports the terraform resources, by default the app fetches the resources only.

`-auto-approve`: Skip interactive approval before importing, requires `-import`

`-terraform_path=path`: Absolute path where to create the terraform folder with all the imported resources. Default is the same folder where you run the tool. Default `./`

`-backend_template=path`: Absolute path of the template file you want to use for the terraform [backend](https://www.terraform.io/docs/backends/config.html) file. Not setting this up along with `-import` will import everything on a local state file, so be aware of that.

## Destination files

Your terraform files will be placed inside a folder called `terraform` with the following structure:

```text
terraform
    modules
        repos
            main.tf
            vatiables.tf
        teams
            main.tf
            vatiables.tf
        users
            main.tf
            vatiables.tf
    users
            backend.tf (If set)
            users.tf
    groups
            backend.tf (If set)
            groups.tf
    repos
            backend.tf (If set)
            repos.tf
```

So you can simply copy the whole content and move it where you need.

I separated every resource to have its state, you normally don't want to manage all resources in one single state, you are going to reach a point where it will be so big that needs ~40/60 mins only to plan (been there, seen that).

### Example of a backend template file for [S3](https://www.terraform.io/docs/backends/types/s3.html)

`backend.tmpl`

```terraform
// remote state on S3
terraform {
  backend "s3" {
    bucket      = "acme/terraform-state/"
    key         = "{{.}}" // This will be replace with the name of the resource the state holds.
    region      = "eu-west-1"
    role_arn    = "arn:aws:iam::00000000000:role/terraform" // remove this if you don't use assume roles
    external_id = "ops-terraform" // same as above
  }
}
```

## Usage | Install

* Set both env vars githubToken and githubOrgANIZATION.
* run the app

## Known issues

* `Users`, `Teams` and `Repos` are the only resources fetched right now, I'll try to add `branch-protection` and `external collaborators`
* It would be also good to have some parameters where you can set the resources you want to import.
* No tests.
