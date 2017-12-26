resource "ns1_record" "it" {
  #required
  zone              = "${ns1_zone.test.zone}"
  domain            = "test.${ns1_zone.test.zone}"
  type              = "CNAME"

  #optional
  ttl               = 60
  use_client_subnet = true
  link = ""

  meta = {
    up          = true
    connections = 5
    latitude    = 0.50
    longitude   = 0.40
  }

  answers = [
    {
      answer = "10.0.0.1"

      meta = {
        up          = true
        connections = 4
        latitude    = 0.5
        georegion   = "US-EAST"
      }
    },
  ]

  regions = [
    {
      name = "cal"

      meta = {
        up          = true
        connections = 3
      }
    },
  ]

  filters {
    filter = "up"
  }

  filters {
    filter = "geotarget_country"
  }

  filters {
    filter = "select_first_n"
    config = {N=1}
  }
}

#records must have an associated zone
resource "ns1_zone" "test" {
  zone = "terraform-record-test.io"
}