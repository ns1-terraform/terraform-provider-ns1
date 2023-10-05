# Geofence Example

resource "ns1_zone" "zone" {
  zone = "geofence.example"
}

resource "ns1_record" "www" {
  zone   = ns1_zone.zone.zone
  domain = "www.${ns1_zone.zone.zone}"
  type   = "A"
  filters {
    filter = "geofence_country"
    config = {
      remove_no_location = 1
    }
  }
  filters {
    filter = "select_first_n"
    config = {
      N = "1"
    }
  }
  answers {
    answer = "1.1.1.1"
    meta = {
      country : "GB",
      note = "UK"
    }
  }
  answers {
    answer = "2.2.2.2"
    meta = {
      country : "CA",
      note = "Canada"
    }
  }
  answers {
    answer = "3.3.3.3"
    meta = {
      note = "Rest of World"
    }
  }
}
