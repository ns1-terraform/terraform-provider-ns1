---
layout: "ns1"
page_title: "NS1: ns1_alert"
sidebar_current: "docs-ns1-resource-alert"
description: |-
  Provides a NS1 Alert resource.
---

# ns1\_alert

Provides a NS1 Alert resource. This can be used to create, modify, and delete alerts.

## Example Usage

```hcl
resource "ns1_alert" "example_zone_alert" {
  #required
  name               = "Example Zone Alert"
  type               = "zone"
  subtype            = "transfer_failed"

  #optional
  notification_lists = []
  zone_names = ["a.b.c.com","myzone"]
  record_ids = []
}

resource "ns1_alert" "example_usage_alert" {
  #required
  name               = "Example Usage Alert"
  type               = "account"
  subtype            = "record_usage"
  data {
    alert_at_percent = 80
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The free-form display name for the alert.
* `type` - (Required) The type of the alert.
* `subtype` - (Required) The type of the alert.
* `notification_lists` - (Optional) A list of id's for notification lists whose notifiers will be triggered by the alert.
* `zone_names` - (Optional) A list of zones this alert applies to.
* `record_ids` - (Optional) A list of record id's this alert applies to.
* `data` - (Optional) A resource block with additional settings: the name and type of them vary based on the alert type.
  * `alert_at_percent` - required by the account/usage alerts, with a value between 1 and 100

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created_by` - (Read Only) The user or apikey that created this alert.
* `updated_by` - (Read Only) The user or apikey that last modified this alert.
* `created_at` - (Read Only) The Unix timestamp representing when the alert configuration was created.
* `updated_at` - (Read Only) The Unix timestamp representing when the alert configuration was last modified.

## Import

`terraform import ns1_alert.<name> <alert_id>`

## NS1 Documentation

[Alerts Api Doc](https://ns1.com/api#alerts)
