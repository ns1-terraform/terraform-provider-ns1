resource "ns1_tsigkey" "example" {
  # Required
  name = "ExampleTsigKey"
  algorithm = "hmac-sha256"
  secret = "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLA=="
}