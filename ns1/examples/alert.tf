resource "ns1_alert" "example" {
  #required
  name               = "Example Alert"
  type               = "zone"
  subtype            = "transfer_failed"

  #optional
  notification_lists = []
  zone_names = []
  record_ids = []
}