# DNS
resource "ns1_monitoringjob" "example_com_dns" {
  job_type        = "dns"
  name            = "[DNS] example.com"
  mute            = false
  active          = true
  regions         = ["nrt", "dal", "sin", "sjc", "lga", "ams", "syd", "gru", "lhr"]
  policy          = "quorum"
  frequency       = 60
  rules {
    key           = "num_records"
    comparison    = "=="
    value         = 1
  }
  rules {
    key           = "rdata"
    comparison    = "contains"
    value         = "93.184.216.34"
  }
  rules {
    key           = "rtt"
    comparison    = "<="
    value         = "100"
  }
  rapid_recheck   = true
  config          = {
    response_timeout = 2000
    domain           = "example.com"
    host             = "8.8.8.8"
    ipv6             = false
    type             = "A"
    port             = 53
  }
  notify_list     = ns1_notifylist.my_notify_list.id
  notify_delay    = 0
  notify_repeat   = 0
  notify_failback = true
  notify_regional = false
  notes           = "example.com A"
}

# HTTP
resource "ns1_monitoringjob" "example_com_http" {
  job_type       = "http"
  name           = "[HTTP] example.com"
  mute           = false
  active         = true
  regions        = ["nrt", "dal", "sin", "sjc", "lga", "ams", "syd", "gru", "lhr"]
  policy         = "quorum"
  frequency      = 60
  rules {
    key        = "status_code"
    comparison = "=="
    value      = "200"
  }
  rules {
    key        = "body"
    comparison = "contains"
    value      = "Example Domain"
  }
  rapid_recheck  = true
  config         = {
    url             = "https://www.example.com/"
    virtual_host    = "example.com"
    method          = "GET"
    user_agent      = "NS1 HTTP Monitoring Job"
    authorization   = "Auth-Token: foobar"
    follow_redirect = false
    connect_timeout = 5
    idle_timeout    = 3
    tls_add_verify  = false
    ipv6            = false
  }
  notify_list     = ns1_notifylist.my_notify_list.id
  notify_delay    = 0
  notify_repeat   = 0
  notify_failback = true
  notify_regional = false
  notes           = "example.com:443 200"
}


# PING
resource "ns1_monitoringjob" "example_com_ping" {
  job_type       = "ping"
  name           = "[PING] example.com"
  mute           = false
  active         = true
  regions        = ["nrt", "dal", "sin", "sjc", "lga", "ams", "syd", "gru", "lhr"]
  policy         = "quorum"
  frequency      = 60
  rules {
    key        = "rtt"
    comparison = "<"
    value      = "300"
  }
  rules {
    key        = "loss"
    comparison = "<="
    value      = "100"
  }
  rapid_recheck  = true
  config         = {
    count    = 4
    host     = "example.com"
    interval = 25
    timeout  = 2000
    ipv6     = false
  }
  notify_list     = ns1_notifylist.my_notify_list.id
  notify_delay    = 0
  notify_repeat   = 0
  notify_failback = true
  notify_regional = false
  notes           = "example.com ICMP"
}

# TCP
resource "ns1_monitoringjob" "example_com_tcp" {
  job_type       = "tcp"
  name           = "[TCP] example.com"
  mute           = false
  active         = true
  regions        = ["nrt", "dal", "sin", "sjc", "lga", "ams", "syd", "gru", "lhr"]
  policy         = "quorum"
  frequency      = 60
  rules {
    key        = "output"
    comparison = "contains"
    value      = "200 OK"
  }
  rules {
    key        = "output"
    comparison = "contains"
    value      = "Example Domain"
  }
  rapid_recheck  = true
  config         = {
    response_timeout = 1000
    ipv6             = false
    send             = "GET / HTTP/1.1\\r\\nHost: example.com\\r\\n\\r\\n"
    connect_timeout  = 2000
    ssl              = 1
    host             = "93.184.216.34"
    tls_add_verify   = false
    port             = 443
  }
  notify_list     = ns1_notifylist.my_notify_list.id
  notify_delay    = 0
  notify_repeat   = 0
  notify_failback = true
  notify_regional = false
  notes           = "example.com GET /"
}
