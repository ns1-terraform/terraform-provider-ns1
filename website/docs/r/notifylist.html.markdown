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
      headers = "Content-Type: application/json"
    }
  }

  notifications {
    type = "email"
    config = {
      email = "test@test.com"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The free-form display name for the notify list.
* `notifications` - (Optional) A list of notifiers. All notifiers in a notification list will receive notifications whenever an event is send to the list (e.g., when a monitoring job fails). Notifiers are documented below.

Notify List Notifiers (`notifications`) support the following:

* `type` - (Required) The type of notifier. Available notifiers are indicated in /notifytypes endpoint.
* `config` - (Required) Configuration details for the given notifier type.
  * `email` - Email to notify to; required for type = "email"
  * `service_key` - Service key of the Pagerduty integration to notify to; required for type = "pagerduty"
  * `sourceid` - Source id of the datafeedto notify to; required for type = "datafeed"
  * `url` - URL to notify to; required for type = "webhook" and "slack"
  * `username` - Username to notify as; required for type = "slack"
  * `channel` - Channel to notify to; required for type = "slack"
  * `headers` - Headers to add in the notification (optional for type = "webhook"): because they're encoded as a string, they have to be in alphabetical order and separated by carriage return, e.g. `"Accept: application/json\nContent-Type: application/json"`

## Attributes Reference

All of the arguments listed above are exported as attributes, with no
additions.

## Import

`terraform import ns1_notifylist.<name> <notifylist_id>`

## NS1 Documentation

[NotifyList Api Doc](https://ns1.com/api#notification-lists)
