resource "ns1_alert" "example1" {
  #required
  name               = "Example Alert"
  type               = "zone"
  subtype            = "transfer_failed"

  #optional
  notification_lists = []
  zone_names = []
  record_ids = []
}

resource "ns1_alert" "example2" {
  #required
  name               = "Example Alert"
	type               = "account"
	subtype            = "record_usage"
	data {
		alert_at_percent = 80
	}
}
