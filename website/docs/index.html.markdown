---
layout: "ns1"
page_title: "Provider: NS1"
sidebar_current: "docs-ns1-index"
description: |-
  The [NS1](https://ns1.com/) provider is used to interact with the resources supported by NS1.
---

# NS1 Provider

The NS1 provider exposes resources to interact with the NS1 REST API. The
provider needs to be configured with the proper credentials before it can be
used. Note also that for a given resource to function, the API key used must
have the corresponding permissions set.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the NS1 provider
provider "ns1" {
  apikey = "${var.ns1_apikey}"
}

# Create a new zone
resource "ns1_zone" "foobar" {
  # ...
}
```

## Argument Reference

The following arguments are supported:

* `apikey` - (Required) NS1 API token. It must be provided, but it can also
  be sourced from the `NS1_APIKEY` environment variable.
* `version` - (Optional, but recommended if you don't like surprises) From
  output of `terraform init`.
* `endpoint` - (Optional) NS1 API endpoint. For managed clients, this normally
  should not be set.
* `ignore_ssl` - (Optional) This normally does not need to be set.
* `rate_limit_parallelism` - (Optional) Integer for parallelism amount
  (terraform's default is 10). If set, uses an alternate strategy to handle
  rate limiting from the NS1 API. Give this a try if you are processing a large
  number of resource and running into `429` rate-limiting errors.

## Environment Variables

The provider does check some environment variables as an alternative to
embedding in the config:

* `NS1_APIKEY` - (string) Explained above.
* `NS1_ENDPOINT` - (string) Explained above.
* `NS1_IGNORE_SSL` - (boolean) If set, follows the convention of
  [strconv.ParseBool](https://golang.org/pkg/strconv/#ParseBool).
* `NS1_RATE_LIMIT_PARALLELISM` - (int) Explained above.
