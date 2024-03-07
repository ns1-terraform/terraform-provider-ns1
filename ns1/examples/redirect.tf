
resource "ns1_redirect_certificate" "example" {
  domain       = "www.example.com"
}

resource "ns1_redirect" "example" {
  domain           = "www.example.com"
  path             = "/from/path"
  target           = "https://url.com/target/path"
  forwarding_mode  = "all"
  forwarding_type  = "permanent"
  https_enabled    = true
  https_forced     = true
  query_forwarding = true
  tags             = []
}
