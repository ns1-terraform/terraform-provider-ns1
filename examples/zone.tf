# Primary Zone
resource "ns1_zone" "primary_example" {
  zone = "terraform.example"
  hostmaster = "hostmaster@terraform.example"
}

# Primary with outgoing XFR
resource "ns1_zone" "primary_example_xfr" {
  zone = "primary.example"
  hostmaster = "hostmaster@terraform.example"
  secondaries {
    ip     = "2.2.2.2"
    port   = 53
    notify = true
  }
  secondaries {
    ip     = "3.3.3.3"
    port   = 5353
    notify = false
  }
}

# Secondary zone
resource "ns1_zone" "secondary_example" {
  zone = "secondary.example"
  primary = "192.0.2.1"
  additional_primaries = ["192.0.2.2"]
}

# Secondary zone with TSIG
resource "ns1_zone" "secondary_example_tsig" {
  zone     = "terraform-tsig.example.io"
  primary  = "1.1.1.1"
  tsig = {
    enabled = true
    name = "terraform_tsig_key"
    hash = "hmac-sha256"
    key = "Ok1qR5IW1ajVka5cHPEJQIXfLyx5V3PSkFBROAzOn21JumDq6nIpoj6H8rfj5Uo+Ok55ZWQ0Wgrf302fDscHLA=="
  }
}
