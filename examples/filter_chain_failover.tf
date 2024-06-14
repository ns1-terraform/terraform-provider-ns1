#Failover Example

resource "ns1_datasource" "monitoring" {
  name       = "monitoring datasource"
  sourcetype = "nsone_monitoring"
}

resource "ns1_notifylist" "failover_notifylist" {
  name = "monitoring data feed notify list"
  notifications {
    type = "datafeed"
    config = {
      sourceid = ns1_datasource.monitoring.id
    }
  }
}

resource "ns1_monitoringjob" "failover_monitor" {
  active = true
  config = {
    "method" = "GET"
    "url"    = "https://www.example.com/"
    "connect_timeout" = "5"
    "idle_timeout" = "3"
    "user_agent" = "NS1 HTTP Monitoring Job"
  }
  frequency = 30
  job_type = "http"
  name = "Terraform Failover monitor"
  notify_failback = true
  notify_list = ns1_notifylist.failover_notifylist.id
  policy = "quorum"
  rapid_recheck = true
  regions = [
          "lhr",
          "dal",
          "gru",
          "sjc",
          "ams",
          "nrt",
          "syd",
  ]
  rules {
    comparison = "=="
    key        = "status_code"
    value      = "200"
  }
}

resource "ns1_datafeed" "feed" {
  name      = "Monitoring datafeed"
  source_id = ns1_datasource.monitoring.id
  config = {
    jobid = ns1_monitoringjob.failover_monitor.id
  }
}

resource "ns1_zone" "zone" {
  zone = "failover.example"
}

resource "ns1_record" "www" {
  zone   = ns1_zone.zone.zone
  domain = "www.${ns1_zone.zone.zone}"
  type   = "A"
  filters {
    filter = "up"
    config = {}
  }
  filters {
    filter = "priority"
    config = {
      eliminate = 1
    }
  }
  filters {
    filter = "select_first_n"
    config = {
      N = "1"
    }
  }
  answers {
    answer = "1.1.1.1"
    meta = {
      up = jsonencode({
        "feed" = ns1_datafeed.feed.id
      }),
      priority = 1,
      note = "Primary answer monitored"
    }
  }
  answers {
    answer = "2.2.2.2"
    meta = {
      up = true,
      priority = 2,
      note = "Failover answer"
    }
  }
}
