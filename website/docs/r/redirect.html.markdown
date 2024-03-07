---
layout: "ns1"
page_title: "NS1: ns1_redirect"
sidebar_current: "docs-ns1-resource-redirect"
description: |-
  Provides a NS1 Redirect resource.
---

# ns1\_redirect

Provides a NS1 Redirect resource. This can be used to create, modify, and delete redirects.

## Example Usage

```hcl
resource "ns1_redirect" "example" {
  domain       = "www.example.com"
  path         = "/from/path"
  target       = "https://url.com/target/path"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) The domain name to redirect from.
* `path` - (Required) The path on the domain to redirect from.
* `target` - (Required) The URL to redirect to.
* `id` - (Optional) The redirect id, if already created.
* `certificate_id` - (Optional) The certificate redirect id, if already created.
* `forwarding_mode` - (Optional - defaults to "all") How the target is interpreted:
  * __all__       appends the entire incoming path to the target destination;
  * __capture__   appends only the part of the incoming path corresponding to the wildcard (*);
  * __none__      does not append any part of the incoming path.
* `forwarding_type` - (Optional - defaults to "permanent") How the redirect is executed:
  * __permanent__ (HTTP 301) indicates to search engines that they should remove the old page from
                  their database and replace it with the new target page (this is recommended for SEO);
  * __temporary__ (HTTP 302) less common, indicates that search engines should keep the old domain or
                  page indexed as the redirect is only temporary (while both pages might appear in the
                  search results, a temporary redirect suggests to the search engine that it should
                  prefer the new target page);
  * __masking__   preserves the redirected domain in the browser's address bar (this lets users see the
                  address they entered, even though the displayed content comes from a different web page).
* `https_enabled` - (Optional - defaults to true) Enables HTTPS support on the source domain by using Let's Encrypt certificates.
* `https_forced` - (Optional - defaults to true) Forces redirect for users that try to visit HTTP domain to HTTPS instead.
* `query_forwarding` - (Optional - defaults to false) Enables the query string of a URL to be applied directly to the new target URL.
* `tags` - (Optional - array) Tags associated with the configuration.
* `last_updated` - (Read Only) The Unix timestamp representing when the redirect configuration was last updated.

## Attributes Reference

All of the arguments listed above are exported as attributes, with no
additions.

## NS1 Documentation

[Redirect Api Doc](https://ns1.com/api#redirect)


# ns1\_redirect\_certificate

Provides a NS1 Redirect Certificate resource. This can be used to create, modify, and delete redirect certificates.

## Example Usage

```hcl
resource "ns1_redirect_certificate" "example" {
  domain       = "www.example.com"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) The domain the redirect refers to.
* `id` - (Optional) The certificate id, if already created.
* `certificate` - (Read Only) The certificate value.
* `valid_from` - (Read Only) The Unix timestamp representing when the certificate first started being valid.
* `valid_until` - (Read Only) The Unix timestamp representing when the certificate will stop being valid.
* `processing` - (Read Only) Whether the certificate is active.
* `errors` - (Read Only) Any error encountered when applying the certificate.
* `last_updated` - (Read Only) The Unix timestamp representing when the certificate was last signed.

## Attributes Reference

All of the arguments listed above are exported as attributes, with no
additions.

## NS1 Documentation

[Redirect Api Doc](https://ns1.com/api#redirect)
