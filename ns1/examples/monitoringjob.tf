resource "ns1_monitoringjob" "it" {
  #required
  job_type = "tcp"
  name     = "terraform test"

  regions   = ["lga","sjc","sin"]
  frequency = 60

  config = {
    ssl  = "1",
    send = "HEAD / HTTP/1.0\\r\\n\\r\\n"
    port = 443
    host = "1.2.3.4"
    connect_timeout = 2000
    ipv6 = false
    response_timeout = 1000
    tls_add_verify = false
  }

  #optional
  active          = true
  rapid_recheck   = false
  notes           = "some notes about this job"
  notify_delay    = 3000
  notify_repeat   = 3000
  notify_failback = false
  notify_list     = ""
  notify_regional = true


  rules {
    value      = "200 OK"
    comparison = "contains"
    key        = "output"
  }
}
