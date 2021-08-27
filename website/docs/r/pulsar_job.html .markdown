---
layout: "ns1"
page_title: "NS1: pulsarjob"
sidebar_current: "docs-ns1-resource-pulsarjob"
description: |-
  Provides a Pulsar Job resource.
---

# ns1\_pulsar\_job

Provides a Pulsar job resource. This can be used to create, modify, and delete zones.

## Example Usage

```hcl
# Create a new Pulsar Application
resource "ns1_application" "example" {
  name = "terraform.example.io"
}

# Create a new Pulsar JavaScript Job with Blend Metric Weights and multiple weights
resource "ns1_pulsarjob" "example_javascript" {
  name    = "terraform.example_javascript.io"
  appid = "${ns1_pulsar_application.example.id}"
  typeid  = "latency"
  
  config = {
    host = "terraform.job_host.io"
    url_path = "/terraform.job_url_path.io"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Pulsar job. Typically, this is the name of the CDN or endpoint.
* `appid` - (Required) ID of the Pulsar app.
* `typeid` - (Required) Specifies the type of Pulsar job - either latency or custom.
* `active` - (Optional) The job's status, if it's active or not.
* `shared` - (Optional) Enable to share data with other approved accounts.
* `Config` - (Optional) [Config](#config-1) is documented below. Note: **Required if typeid is "latency"** 


#### Config

`config` support several sub-options related to job configuration:

* `host` - (Required) Hostname where the resource is located.
* `url_path` - (Required) URL path to be appended to the host.
* `https` - (Optional) Indicates whether or not to use HTTPS in measurements.
* `http` - (Optional) Indicates whether or not to use HTTP in measurements.
* `request_timeout_millis` - (Optional) The amount of time to allow a single job to perform N runs.
* `job_timeout_millis` - (Optional) The amount of time to allow a single job to perform 1 run.
* `use_xhr` - (Optional) Indicates wheter or not to use XmlHttpRequest (XHR) when taking measurements.
* `satuc_values` - (Optional) Indicates wheter or not to skip aggregation for this job's measurements.

## Import

`terraform import ns1_pulsarjob.<name> <appid>_<jobid>`

## NS1 Documentation

[Pulsar Job Api Docs](https://ns1.com/api#jobs)
