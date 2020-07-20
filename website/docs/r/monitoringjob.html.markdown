---
layout: "ns1"
page_title: "NS1: ns1_monitoringjob"
sidebar_current: "docs-ns1-resource-monitoringjob"
description: |-
  Provides a NS1 Monitoring Job resource.
---

# ns1\_monitoringjob

Provides a NS1 Monitoring Job resource. This can be used to create, modify, and delete monitoring jobs.

## Example Usage

```hcl
resource "ns1_monitoringjob" "uswest_monitor" {
  name          = "uswest"
  active        = true
  regions       = ["sjc", "sin", "lga"]
  job_type      = "tcp"
  frequency     = 60
  rapid_recheck = true
  policy        = "quorum"

  config = {
    ssl  = 1
    send = "HEAD / HTTP/1.0\r\n\r\n"
    port = 443
    host = "example-elb-uswest.aws.amazon.com"
  }

  rules {
    value      = "200 OK"
    comparison = "contains"
    key        = "output"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The free-form display name for the monitoring job.
* `job_type` - (Required) The type of monitoring job to be run. Refer to the NS1 API documentation (https://ns1.com/api#monitoring-jobs) for supported values which include ping, tcp, dns, http.
* `active` - (Required) Indicates if the job is active or temporarily disabled.
* `regions` - (Required) The list of region codes in which to run the monitoring
  job. See NS1 API docs for supported values.
* `frequency` - (Required) The frequency, in seconds, at which to run the monitoring job in each region.
* `rapid_recheck` - (Required) If true, on any apparent state change, the job is quickly re-run after one second to confirm the state change before notification.
* `policy` - (Required) The policy for determining the monitor's global status
  based on the status of the job in all regions. See NS1 API docs for supported values.
* `config` - (Required) A configuration dictionary with keys and values depending on the job_type. Configuration details for each job_type are found by submitting a GET request to https://api.nsone.net/v1/monitoring/jobtypes.
* `notify_delay` - (Optional) The time in seconds after a failure to wait before sending a notification.
* `notify_repeat` - (Optional) The time in seconds between repeat notifications of a failed job.
* `notify_failback` - (Optional) If true, a notification is sent when a job returns to an "up" state.
* `notify_regional` - (Optional) If true, notifications are sent for any regional failure (and failback if desired), in addition to global state notifications.
* `notify_list` - (Optional) The Terraform ID (e.g. ns1_notifylist.my_slack_notifier.id) of the notification list to which monitoring notifications should be sent.
* `notes` - (Optional) Freeform notes to be included in any notifications about this job.
* `rules` - (Optional) A list of rules for determining failure conditions. Each rule acts on one of the outputs from the monitoring job. You must specify key (the output key); comparison (a comparison to perform on the the output); and value (the value to compare to). For example, {"key":"rtt", "comparison":"<", "value":100} is a rule requiring the rtt from a job to be under 100ms, or the job will be marked failed. Available output keys, comparators, and value types are are found by submitting a GET request to https://api.nsone.net/v1/monitoring/jobtypes.

## Attributes Reference

All of the arguments listed above are exported as attributes, with no
additions.

## NS1 Documentation

[MonitoringJob Api Doc](https://ns1.com/api#monitoring-jobs)
