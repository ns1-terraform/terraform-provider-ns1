---
layout: "ns1"
page_title: "NS1: ns1_notifylist"
sidebar_current: "docs-ns1-resource-notifylist"
description: |-
  Provides a NS1 Notify List resource.
---

# ns1\_notifylist

Provides a NS1 Notify List resource. This can be used to create, modify, and delete notify lists.

## Example Usage

```hcl
resource "ns1_notifylist" "nl" {
  name = "my notify list"
  notifications {
    type = "webhook"
    config = {
      url = "http://www.mywebhook.com"
    }
  }

  notifications {
    type = "email"
    config = {
      email = "test@test.com"
    }
  }
  
 notifications {
    config = {
      sourceid = "${ns1_datasource.monitoring_datasource.id}"
    }
    type = "datafeed"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The free-form display name for the notify list.
* `notifications` - (Optional) A list of notifiers. All notifiers in a notification list will receive notifications whenever an event is send to the list (e.g., when a monitoring job fails). Notifiers are documented below.

Notify List Notifiers (`notifications`) support the following:

* `type` - (Required) The type of notifier. Available notifiers are indicated in /notifytypes endpoint. Note type `datafeed` is for NS1 Monitoring Data Sources to enable NS1 monitoring to pass metadata to your records and answers.  
* `config` - (Required) Configuration details for the given notifier type. `sourceid` is required in `config` for notifcations to NS1 Monitoring Data Sources. 

## Attributes Reference

All of the arguments listed above are exported as attributes, with no
additions.

## NS1 Documentation

[NotifyList Api Doc](https://ns1.com/api#notification-lists)
