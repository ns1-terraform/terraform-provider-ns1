---
layout: "ns1"
page_title: "NS1: ns1_record"
sidebar_current: "docs-ns1-resource-record"
description: |-
  Provides a NS1 Record resource.
---

# ns1\_record

Provides a NS1 Record resource. This can be used to create, modify, and delete records.

## Example Usage

```hcl
resource "ns1_zone" "example" {
  zone = "terraform.example.io"
}

resource "ns1_datasource" "ns1" {
  name       = "ns1_source"
  sourcetype = "nsone_v1"
}

resource "ns1_datafeed" "foo" {
  name      = "foo_feed"
  source_id = "${ns1_datasource.ns1.id}"
  config = {
    label = "foo"
  }
}

resource "ns1_datafeed" "bar" {
  name      = "bar_feed"
  source_id = "${ns1_datasource.ns1.id}"
  config = {
    label = "bar"
  }
}

resource "ns1_record" "www" {
  zone   = "${ns1_zone.tld.zone}"
  domain = "www.${ns1_zone.tld.zone}"
  type   = "CNAME"
  ttl    = 60
  meta   = {
    up   = true
  }

  regions {
    name = "east"
    meta = {
      georegion = "US-EAST"
    }
  }

  regions {
    name = "usa"
    meta = {
      country = "US"
    }
  }

  answers {
    answer  = "sub1.${ns1_zone.tld.zone}"
    region  = "east"
    meta    = {
      up    = "{\"feed\":\"${ns1_datafeed.foo.id}\"}"
    }
  }

  answers {
    answer = "sub2.${ns1_zone.tld.zone}"
    meta   = {
      up   = "{\"feed\":\"${ns1_datafeed.bar.id}\"}"
      connections = 3
    }
  }

  filters {
    filter = "select_first_n"

    config = {
      N = 1
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone` - (Required) The zone the record belongs to.
* `domain` - (Required) The records' domain.
* `type` - (Required) The records' RR type.
* `ttl` - (Optional) The records' time to live.
* `link` - (Optional) The target record to link to. This means this record is a
  'linked' record, and it inherits all properties from its target.
* `use_client_subnet` - (Optional) Whether to use EDNS client subnet data when
  available(in filter chain).
* ` meta` - (Optional) meta is supported at the `record` level. [Meta](#meta-3)
  is documented below.
* `regions` - (Optional) One or more "regions" for the record. These are really
  just groupings based on metadata, and are called "Answer Groups" in the NS1 UI,
  but remain `regions` here for legacy reasons. [Regions](#regions-1) are
  documented below. Please note the ordering requirement!
* `answers` - (Optional) One or more NS1 answers for the records' specified type.
  [Answers](#answers-1) are documented below.
* `filters` - (Optional) One or more NS1 filters for the record(order matters).
  [Filters](#filters-1) are documented below.

#### Answers

`answers` support the following:

* `answer` - (Required) Space delimited string of RDATA fields dependent on the record type.

    A:

        answer = "1.2.3.4"

    CNAME:

        answer = "www.example.com"

    MX:

        answer = "5 mail.example.com"

    SRV:

        answer = "10 0 2380 node-1.example.com"

    SPF:

        answer = "v=DKIM1; k=rsa; p=XXXXXXXX"

   
* `region` - (Optional) The region (Answer Group really) that this answer
  belongs to. This should be one of the names specified in `regions`. Only a
  single `region` per answer is currently supported. If you want an answer in
  multiple regions, duplicating the answer (including metadata) is the correct
  approach.
* ` meta` - (Optional) meta is supported at the `answer` level. [Meta](#meta-3)
  is documented below.

#### Filters

`filters` support the following:

* `filter` - (Required) The type of filter.
* `disabled` - (Optional) Determines whether the filter is applied in the
  filter chain.
* `config` - (Optional) The filters' configuration. Simple key/value pairs
  determined by the filter type.

#### Regions

`regions` support the following:

* `name` - (Required) Name of the region (or Answer Group).
* `meta` - (Optional) meta is supported at the `regions` level. [Meta](#meta-3)
  is documented below.
  Note that `Meta` values for `country`, `ca_province`, `georegion`, and
  `us_state` should be comma separated strings, and changes in ordering will not
  lead to terraform detecting a change.

Note: regions **must** be sorted lexically by their "name" argument in the
Terraform configuration file, otherwise Terraform will detect changes to the
record when none actually exist.

#### Meta

Metadata (`meta`) is a bit tricky at the moment. For "static" values it works
as you would expect, but when a value is a `datafeed`, it should be represented
as "escaped" JSON. See the [Example Usage](#example-usage) above for
illustration of this.

Note that variables are still supported in the escaped JSON format. Note also
that we intend to fix up this "escaped" JSON stuff as soon as possible, so
please bear with us and plan accordingly.

Since this resource supports [import](#import), you may find it helpful to set
up some `meta` fields via the web portal or API, and use the results from
import to ensure that everything is properly escaped and evaluated.

See [NS1 API](https://ns1.com/api#get-available-metadata-fields) for the most
up-to-date list of available `meta` fields.

## Attributes Reference

All of the arguments listed above are exported as attributes, with no
additions.

## Import

`terraform import ns1_record.<name> <zone>/<domain>/<type>`

So for the example above:

`terraform import ns1_record.www terraform.example.io/www.terraform.example.io/CNAME`
