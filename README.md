NS1 Terraform Provider
==================

> This project is in [active development](https://github.com/ns1/community/blob/master/project_status/ACTIVE_DEVELOPMENT.md).

- NS1 Website: https://www.ns1.com
- Terraform Website: https://www.terraform.io
- Terraform NS1 Provider Documentation: https://registry.terraform.io/providers/ns1-terraform/ns1/latest/docs
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Contents
------
1. [Upgrading from Terraform 0.12](#upgrading-from-terraform-012) - considerations when upgrading from previous versions of Terraform
2. [Requirements](#requirements) - lists the requirements for building the provider
3. [Building The Provider](#building-the-provider) - lists the steps for building the provider
4. [Using The Provider](#using-the-provider) - details how to use the provider
5. [Developing The Provider](#developing-the-provider) - steps for contributing back to the provider
6. [Known Isssues/Roadmap](#known-issuesroadmap) - check here for some of the improvements we are working on

Upgrading from Terraform 0.12
-----------------------------
In preperation for the 0.13 release of Terraform, this repo has recently changed locations from the Hashicorp GitHub org to one owned by NS1.

As a result, in order to upgrade an existing config and state to Terraform 0.13 and use NS1 provider v1.8.5 and above, you'll need to
update your config's `required_providers` block to point to the new location. 

The `0.13upgrade` tool will display a warning suggesting as much, but will not enforce this or automatically update your config.

Note that this block is only required if updating an existing state from `0.12` or below.  Fresh deployments via Terraform `0.13` or above will automatically detect the new location of the provider.

Here is an example `required_providers` block enforcing the new location of this provider and Terraform `0.13` or greater:
```
terraform {
  required_providers {
    ns1 = {
      source = "ns1-terraform/ns1"
    }
  }
  required_version = ">= 0.13"
}
```

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.13+
-	[Go](https://golang.org/doc/install) 1.12+ (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/ns1-terraform/terraform-provider-ns1`

```sh
$ mkdir -p $GOPATH/src/github.com/ns1-terraform
$ cd $GOPATH/src/github.com/ns1-terraform
$ git clone git@github.com:ns1-terraform/terraform-provider-ns1.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/ns1-terraform/terraform-provider-ns1
$ make build
```

Using The Provider
----------------------

The documentation and examples for NS1 `Resources` and `Data Sources` is
maintained as part of this repository, in the `/website` directory. This is
published to
[registry.terraform.io/providers/ns1-terraform/ns1/latest/docs](https://registry.terraform.io/providers/ns1-terraform/ns1/latest/docs)
as part of the release process.


Developing The Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine 
(version 1.12+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH),
as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in 
the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-ns1
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

Some helpful things for debugging:

* Set `TF_LOG=DEBUG` for verbose logging.
* Additionally set `NS1_DEBUG` environment variable to include details of the
  API requests in the logs.

Contributions
---

Pull Requests and issues are welcome. See the [NS1 Contribution Guidelines](https://github.com/ns1/community) for more information.

Known Issues/Roadmap
--------------------

* Currently, some arguments marked as required in resource documentation are
  de-facto optional. A resource will be created/updated without error, but
  in general will lead to a "dirty terraform" state, since the defaulted
  attributes on the returned state may not match the resource descriptions.
  We're working on making these either truly Required or truly Optional as
  appropriate.
* Currently, some resources do not return attributes for optional features that
  are unused. We are working on making the resource schemas fixed, with proper
  defaults returned for optional/unused features.
