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

resource "ns1_alert_sso" "example_saml_alert" {
  #required
  name               = "Example Alert"
  type               = "account"
  subtype            = "saml_certificate_expired"
  notification_lists = []
} 

resource "ns1_alert_redirect" "example_redirect_alert" {
  #required
  name               = "Example Alert"
  type               = "redirects"
  subtype            = "certificate_renewal_failed"
  notification_lists = []
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