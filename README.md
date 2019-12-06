NS1 Terraform Provider
==================

- NS1 Website: https://www.ns1.com
- Terraform Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Contents
------
1. [Requirements](#requirements) - lists the requirements for building the provider
2. [Building The Provider](#building-the-provider) - lists the steps for building the provider
3. [Using The Provider](#using-the-provider) - details how to use the provider
4. [Developing The Provider](#developing-the-provider) - steps for contributing back to the provider
5. [Known Isssues/Roadmap](#known-issues) - check here for some of the improvements we are working on

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.12+ (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-ns1`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-ns1
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-ns1
$ make build
```

Using The Provider
----------------------

### NS1 Resources List and examples

1. [ApiKey](#apikey)
2. [Datafeed](#datafeed)
3. [Datasource](#datasource)
4. [MonitoringJob](#monitoringjob)
5. [NotifyList](#notifylist)
6. [Record](#record)
7. [Team](#team)
8. [User](#user)
9. [Zone](#zone)

### Addendum

1. [Permissions](#permissions)

### ApiKey

[ApiKeys Api Doc](https://ns1.com/api#api-key)

ApiKeys are one of the data types that supports permissions at the top level of its Terraform resource
in addition to its regular parameters. See [Permissions](#permissions) for the parameters that are available.

_Example_

```hcl
resource "ns1_apikey" "apikey" {
  #required
  name = "my api key"

  #optional
  teams = ["myteam"]

  #permissions are available at the top level, see permissions for documentation
}
```

### Datafeed

[Datafeed Api Doc](https://ns1.com/api#data-feeds)

A Datafeed _requires_ a [Datasource](#datasource)

_Example_
```hcl
resource "ns1_datafeed" "datafeed" {
  name = "terraform test"
  source_id = "${ns1_datasource.api.id}"

  #optional
  config {
    label = "exampledc2"
  }
}
```

### Datasource

[Datasource Api Doc](https://ns1.com/api#data-sources)

_Example_
```hcl
resource "ns1_datasource" "datasource" {
  #required
  name = "terraform test"
  sourcetype = "nsone_v1"

  #optional
  config {}
}
```

### MonitoringJob

[MonitoringJob Api Doc](https://ns1.com/api#monitoring-jobs)

_Example_
```hcl
resource "ns1_monitoringjob" "it" {
  #required
  job_type = "tcp"
  name     = "terraform test"

  regions   = ["lga"]
  frequency = 60

  config = {
    ssl = "1",
    send = "HEAD / HTTP/1.0\r\n\r\n"
    port = 443
    host = "1.2.3.4"
  }

  #optional
  active = true
  rapid_recheck = false
  notes = "some notes about this job"
  notify_delay = 3000
  notify_repeat = 3000
  notify_failback = false
  notify_list = ""
  notify_regional = true


  rules = [{
    value = "200 OK"
    comparison = "contains"
    key = "output"
  }]
}
```

### NotifyList

[NotifyList Api Doc](https://ns1.com/api#notification-lists)

_Example_
```hcl
resource "ns1_notifylist" "test" {
  #required
  name = "terraform test"

  #optional
  notifications = {
    type = "webhook"
    config = {
      url = "http://localhost:9090"
    }
  }
}
```

### Record

[Record Api Doc](https://ns1.com/api#records)

Records have metadata at three different levels:

* Record Level - Lowest precedence
* Region Level - Middle precedence
* Answer Level - Highest precedence

Due to some limitations in Terraform's support of nested maplike objects,
there are some irregularities in supporting metadata, however metadata is
now supported at every level. See the documentation for `ns1_record` for
more details and examples.

Note that regions should be sorted by name in the record's regions list,
otherwise terraform will detect changes to the record when none actually exist.

A record _requires_ a [Zone](#zone)

The **zone** and **domain** fields should not have any leading or trailing dots (".").
If the value is coming from another resource with a leading or trailing dot, it should be cleaned:

`zone = replace(".terraform-test-zone.io.", "/(^\\.)|(\\.$)/", "")`

_Example_
```hcl
resource "ns1_record" "it" {
  #required
  zone              = "${ns1_zone.test.zone}"
  domain            = "test.${ns1_zone.test.zone}"
  type              = "CNAME"

  #optional
  ttl               = 60
  use_client_subnet = true
  link = ""

  meta = {
    up          = true
    connections = 5
    latitude    = 0.50
    longitude   = 0.40
  }

  answers = [
    {
      answer = "10.0.0.1"

      meta = {
        up          = true
        connections = 4
        latitude    = 0.5
        georegion   = "US-EAST"
      }
    },
  ]

  regions = [
    {
      name = "cal"

      meta = {
        up          = true
        connections = 3
      }
    },
  ]

  filters {
    filter = "up"
  }

  filters {
    filter = "geotarget_country"
  }

  filters {
    filter = "select_first_n"
    config = {N=1}
  }
}
```

### Team

[Team Api Docs](https://ns1.com/api#team)

Team is one of the data types that supports permissions at the top level of its Terraform resource
in addition to its regular parameters. See [Permissions](#permissions) for the parameters that are available.


_Example_
```hcl
resource "ns1_team" "foobar" {
  name = "terraform test"

  dns_view_zones = true
  dns_zones_allow_by_default = true
  dns_zones_allow = ["mytest.zone"]
  dns_zones_deny = ["myother.zone"]

  data_manage_datasources = true
}
```

### User

[User Api Docs](https://ns1.com/api#user)

User is one of the data types that supports permissions at the top level of its Terraform resource
in addition to its regular parameters. See [Permissions](#permissions) for the parameters that are available.

_Example_
```hcl
resource "ns1_user" "u" {
  #required
  name = "terraform acc test user %s"
  username = "tf_acc_test_user_%s"
  email = "tf_acc_test_ns1@hashicorp.com"

  #optional
  teams = ["${ns1_team.t.id}"]
  notify {
    billing = true
  }
}
```

### Zone

[Zone Api Docs](https://ns1.com/api#zones)

_Example_
```hcl
resource "ns1_zone" "it" {
  zone    = "terraform-test-zone.io"
  ttl     = 10800
  refresh = 3600
  retry   = 300
  expiry  = 2592000
  nx_ttl  = 3601
}
```

### Permissions

There are three resources that support permissions:

* Apikey
* Team
* User

For each of those resources, these parameters are available at the top level of the resource:

* dns_view_zones: boolean - allows the requestor to view zones
* dns_manage_zones: boolean - allows the requestor to edit/manage zones
* dns_zones_allow_by_default: boolean
* dns_zones_deny: list of strings - explicitly deny these zones for this user/team/key
* dns_zones_allow: list of strings - explicitly allow these zones for this user/team/key
* data_push_to_datafeeds: boolean - allows the requestor to push to datafeeds
* data_manage_datasources: boolean - allows the requestor to manage datasources
* data_manage_datafeeds: boolean - allows the requestor to manage datafeeds
* account_manage_users: boolean - allows the requstor to manage users
* account_manage_payment_methods: boolean - allows the requestor to manage payment methods
* account_manage_plan: boolean - allows the requestor to manage the account payment plan
* account_manage_teams: boolean - allows the requestor to manage teams
* account_manage_apikeys: boolean - allows the requestor to manage apikeys
* account_manage_account_settings: boolean - allows the requestor to manage account settings
* account_view_activity_log: boolean - allows the requestor to view the activity log
* account_view_invoices: boolean - allows the requestor to view account invoices
* monitoring_manage_lists: boolean - allows the requestor to manage monitoring lists
* monitoring_manage_jobs: boolean - allows the requestor to manage monitoring jobs
* monitoring_view_jobs: boolean - allows the requestor to view monitoring jobs

Developing The Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine 
(version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH),
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
* We'll be adding a `record` data source ASAP, to cover simple read-only use
  cases
