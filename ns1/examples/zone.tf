resource "ns1_zone" "it" {
  zone    = "terraform-test-zone.io"
  ttl     = 10800
  refresh = 3600
  retry   = 300
  expiry  = 2592000
  nx_ttl  = 3601
}