resource "ns1_notifylist" "my_notify_list" {
## Notify to a Monitoring Data Source
  name = "My Notify List"
  notifications {
    type = "datafeed"
    config = {
      sourceid = ns1_datasource.my_monitoring_data_source.id
    }
  }
## Notify to an e-mail address
#  notifications {
#    type = "email"
#    config = {
#      email = "jdoe@example.com"
#    }
#  }
## Notify to a Slack Channel
#  notifications {
#    type = "slack"
#    config = {
#      url      = "https://example.slack.com/services/.../.../..."
#      channel  = "#channeltonotify"
#      username = "nametodisplayas"
#    }
#  }
## Notify to a Slack User
#  notifications {
#    type = "slack"
#    config = {
#      url = "https://example.slack.com/services/.../.../..."
#      user = "usernametosendto"
#    }
#  }
## Notify to a Web Hook
#  notifications {
#    type = "webhook"
#    config = {
#      url = "http://www.example.com"
#    }
#  }
## Notify to PagerDuty
#  notifications {
#    type = "pagerduty"
#    config = {
#      service_key = "PAGER_DUTY_SERVICE_KEY"
#    }
#  }
}
