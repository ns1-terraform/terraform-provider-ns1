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

resource "ns1_alert_sso" "example" {
  #required
  name               = "Example Alert"
  type               = "account"
  subtype            = "saml_certificate_expired"
  notification_lists = []
} 

resource "ns1_alert_redirect" "example" {
  #required
  name               = "Example Alert"
  type               = "redirects"
  subtype            = "certificate_renewal_failed"
  notification_lists = []
} 