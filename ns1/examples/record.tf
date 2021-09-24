resource "ns1_record" "it" {
  #required
  zone   = ns1_zone.test.zone
  domain = "test.${ns1_zone.test.zone}"
  type   = "CNAME"

  #optional
  ttl               = 60
  use_client_subnet = true
  link              = ""

  meta = {
    up          = true
    connections = 5
    latitude    = 0.50
    longitude   = 0.40
  }

  answers {
    answer = "a.example.com."
    region = "ExampleRegionA"

    meta = {
      up          = true
      connections = 4
      latitude    = 0.5
      georegion   = "US-EAST"
    }
  }

  answers {
    answer = "b.example.com."
    region = "ExampleRegionB"
  }

  regions {
    name = "ExampleRegionA"

    meta = {
      up          = true
      connections = 3
      country     = "A3,PM"
      us_state    = "AK,AP,AR,AZ,CA,CO,HI,ID,KS,MT,ND,NE,NV,OR,SD,U1,UM,UT,WA,WY"
      ca_province = "AB,BC,NT,SK,U2,YT"
    }
  }

  regions {
    name = "ExampleRegionB"
    meta = {
      country = "AS,AU,CC,CK,CX,FJ,FM,GU,HM,KI,MH,MP,NC,NF,NR,NU,NZ,PF,PG,PN,PW,SB,TK,TO,TV,U9,VU,WF,WS"
      up      = false
    }
  }

  filters {
    filter = "up"
  }

  filters {
    filter = "geotarget_country"
  }

  filters {
    filter = "select_first_n"
    config = { N = 1 }
  }
}

#records must have an associated zone
resource "ns1_zone" "test" {
  zone = "terraform-record-test.io"
}
