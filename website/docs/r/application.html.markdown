---
layout: "ns1"
page_title: "NS1: ns1_application"
sidebar_current: "docs-ns1-resource-application"
description: |- Provides a NS1 Pulsar Application resource.
---

# ns1_application

Provides a NS1 Pulsar application resource. This can be used to create, modify, and delete applications.

## Example Usage

```hcl
# Create new basic pulsar application
resource "ns1_application" "ns1_app" {
  name = "terraform_app"
}

# Create new pulsar application
resource "ns1_application" "ns1_app" {
  name = "terraform_app"
  active = true
  browser_wait_millis = 100
  jobs_per_transaction = 100
}

# Create a new pulsar application with default config
resource "ns1_application" "ns1_app" {
  name = "terraform_app"
  default_config {
    http     = true
    https = false
    request_timeout_millis = 100
    job_timeout_millis = 100
    static_values = true
  }
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Descriptive name for this Pulsar app.
* `active` - (Optional)    Indicates whether or not this application is currently active and usable for traffic
  steering.
* `browser_wait_millis` - (Optional) The amount of time (in milliseconds) the browser should wait before running
  measurements.
* `jobs_per_transaction` -(Optional) Number of jobs to measure per user impression.
* `default_config` -(Optional) Default job configuration. If a field is present here and not on a specific job
  associated with this application, the default value specified here is used..

#### Default config

* `http` - (Optional) Indicates whether or not to use HTTP in measurements.
* `https` - (Optional) Indicates whether or not to use HTTPS in measurements.
* `request_timeout_millis` - (Optional) Maximum timeout per request.
* `job_timeout_millis` - (Optional) - Maximum timeout per job
  0, the primary NSONE Global Network. Normally, you should not have to worry about this.
* `use_xhr` - (Optional) - Whether to use XMLHttpRequest (XHR) when taking measurements.
* `static_values` - (Optional) - Indicates whether or not to skip aggregation for this job's measurements

## Import

`terraform import ns1_application`

So for the example above:

`terraform import ns1_application.example terraform.example.io`

## NS1 Documentation

[Application Api Docs](https://ns1.com/api#get-list-pulsar-applications)
