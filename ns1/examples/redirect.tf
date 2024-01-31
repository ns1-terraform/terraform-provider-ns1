
resource "ns1_redirect_certificate" "example" {
  domain       = "www.example.com"
}

resource "ns1_redirect" "example" {
  domain           = "www.example.com"
  path             = "/from/path"
  target           = "https://url.com/target/path"
  forwarding_mode  = "all"
  forwarding_type  = "permanent"
  ssl_enabled      = true
  force_redirect   = true
  query_forwarding = true
  tags             = []
}
