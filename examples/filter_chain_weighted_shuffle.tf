# Weighted Shuffle Example

resource "ns1_zone" "zone" {
  zone = "shuffle.example"
}

resource "ns1_record" "www" {
  zone   = ns1_zone.zone.zone
  domain = "www.${ns1_zone.zone.zone}"
  type   = "A"
  filters {
    filter = "weighted_shuffle"
    config = {}
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
      weight = 50
    }
  }
  answers {
    answer = "2.2.2.2"
    meta = {
      weight = 50
    }
  }
}
