---
layout: "ns1"
page_title: "NS1: ns1_user_team_member"
sidebar_current: "docs-ns1-resource-user-team-member"
description: |-
  Provides a NS1 User resource.
---

# ns1\_user\_team\_member

Provides a NS1 User resource. Creating a user sends an invitation email to the user's email address. This can be used to create, modify, and delete users.

This resource is specifically to manage a user that is a member of a team.  It is identical to the `ns1_user` resource with the exception
of permission attributes, which are inherited from their team.

## Example Usage

```hcl
resource "ns1_team" "example" {
  name = "Example team"

  dns_view_zones       = false
  account_manage_users = false
}

resource "ns1_user_team_member" "example" {
  name     = "Example User"
  username = "example_user"
  email    = "user@example.com"
  teams    = ["${ns1_team.example.id}"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The free form name of the user.
* `username` - (Required) The users login name.
* `email` - (Required) The email address of the user.
* `notify` - (Required) Whether or not to notify the user of specified events. Only `billing` is available currently.
* `teams` - (Required) The teams that the user belongs to.

