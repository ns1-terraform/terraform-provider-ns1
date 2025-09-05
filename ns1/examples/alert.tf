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
