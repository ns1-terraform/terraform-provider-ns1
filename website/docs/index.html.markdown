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
terraform {
  required_providers {
    ns1 = {
      source = "ns1-terraform/ns1"
    }
  }
  required_version = ">= 0.13"
}

provider "ns1" {
  apikey = var.ns1_apikey
  rate_limit_parallelism = 60
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
* `version` - (Optional, but recommended if you don't like surprises) Require a specific version of the NS1 provider. Run `terraform init` to get your current version.
* `retry_max` - (Optional, introduced in v1.13.2) Sets the number of retries for 50x-series errors. The default is 3. A negative value such as -1 disables this feature.
* `endpoint` - (Optional) NS1 API endpoint. Normally not set unless testing or using non-standard proxies.
* `ignore_ssl` - (Optional) This normally does not need to be set.
* `enable_ddi` - (Deprecated) Enable the DDI-compatible permissions schema. No longer in use.
* `rate_limit_parallelism` - (Optional) Integer for alternative rate limit and parallelism strategy.
    NS1 uses a token-based method for rate limiting API requests. Full details can be found at https://help.ns1.com/hc/en-us/articles/360020250573-About-API-rate-limiting.
    
    By default, the NS1 provider uses the "sleep" strategy of the underlying [NS1 Go SDK](https://github.com/ns1/ns1-go) for handling the NS1 API rate limit:
    an operation waits after every API request for a time equal to the rate limit period of that request type divided by the corresponding tokens remaining.
    
    Furthermore, the default behaviour of Terraform uses ten concurrent operations.
    This means that the provider will burst through available request tokens, gradually slowing until it reaches an equilibrium point where the ten operations wait long enough between requests to replenish ten tokens.
    However, if there are other concurrent uses of the API this can lead to the tokens being entirely depleted when a Terraform operation makes a new request.
    This results in a 429 response which will cause the entire Terraform process to fail.
    
    If you encounter this scenario, or believe you are likely to, then you can set the `rate_limit_parallelism` to enable an alternative rate limiting strategy.
    Here the Terraform operations will burst through all available tokens until they reach a point where the remaining limit is less, or equal, to the value set;
    after this point an operation will wait for the time it would take to replenish an equal number of tokens.
    
    Setting this to a value of 60 represents a good balance between optimising for performance and reducing the risk of a 429 response.
    If you still encounter issues then you can increase this value: we would recommend you do so in increments of 20.
    
    Note: We recommend that you NOT set the Terraform command line `-parallelism=n` option when you run `terraform apply`.
    The default value of ten is sufficient - increasing it will lead to a greater risk of encountering a 429 response.

## Environment Variables

The provider honors the following environment variables for its configuration
variables if they are not specified in the configuration:

* `NS1_APIKEY` - (string) Explained above.
* `NS1_ENDPOINT` - (string) Explained above.
* `NS1_RETRY_MAX` - (integer) Explained above.
* `NS1_IGNORE_SSL` - (boolean) If set, follows the convention of
  [strconv.ParseBool](https://golang.org/pkg/strconv/#ParseBool).
* `NS1_RATE_LIMIT_PARALLELISM` - (integer) Explained above.
* `NS1_TF_USER_AGENT` - (string) Sets the User-Agent header in the NS1 API.
