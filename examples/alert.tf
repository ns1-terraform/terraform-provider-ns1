resource "ns1_alert" "email_on_zone_transfer_failure" {
  name               = "Zone transfer failed"
  type               = "zone"
  subtype            = "transfer_failed"
  notification_lists = [ns1_notifylist.email_list.id]
  zone_names = [ns1_zone.alert_example_one.zone, ns1_zone.alert_example_two.zone]
}

# Nofitication list
resource "ns1_notifylist" "email_list" {
  name = "email list"
  notifications {
    type = "email"
    config = {
      email = "jdoe@example.com"
    }
  }
}

# Secondary zones
resource "ns1_zone" "alert_example_one" {
  zone                 = "alert1.example"
  primary              = "192.0.2.1"
  additional_primaries = ["192.0.2.2"]
}

resource "ns1_zone" "alert_example_two" {
  zone                 = "alert2.example"
  primary              = "192.0.2.1"
  additional_primaries = ["192.0.2.2"]
}