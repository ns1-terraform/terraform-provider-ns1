## Record Examples

# Zone needed
resource "ns1_zone" "zone_example" {
  zone = "terraform.example"
  hostmaster = "hostmaster@terraform.example"
}

# A record
resource "ns1_record" "a_record" {
  zone   = ns1_zone.zone_example.zone
  domain = "a.${ns1_zone.zone_example.zone}"
  type   = "A"
}

# CNAME record resource "ns1_record" "a_record" {
resource "ns1_record" "cname__record" {
  zone   = ns1_zone.zone_example.zone
  domain = "cname.${ns1_zone.zone_example.zone}"
  type   = "CNAME"
  answers {
    answer = "www.example.com"
  }
}

# URLFWD (HTTP redirect)
resource "ns1_record" "url_forward_record" {
  zone   = ns1_zone.zone_example.zone
  domain = "redirect.${ns1_zone.zone_example.zone}"
  type   = "URLFWD"
  answers {
    answer = "/ https://ns1.com 301 2 0"
  }
}

# ALIAS
resource "ns1_record" "alias_record" {
  zone   = ns1_zone.zone_example.zone
  domain = ns1_zone.zone_example.zone
  type   = "ALIAS"
  answers {
    answer = "www.example.com"
  }
}
