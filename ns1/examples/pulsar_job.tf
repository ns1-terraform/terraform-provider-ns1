resource "ns1_pulsarjob" "example" {
  # Required
  name = "<Pulsar job name>"
  appid = "<your Application ID>"
  typeid = "custom"
  
  /* Obs:
    If typeid is "latency", host and url_path becomes required.

    If blend_metric_weights is setted, host, url_path, timestamp 
    and weights becomes required.
  */

  # Optional
  active = false
  shared = false
  config = {
    host = "my_host.com"
    url_path = "/my_url_path"
    https = false
    http = false
    request_timeout_millis = 123
    job_timeout_millis = 321
    use_xhr = false
    static_values = false
  }
  blend_metric_weights = {
    timestamp = 53
  }
  weights {
      name = "WeightName1"
      weight = 3
      default_value = 5.2
      maximize = false
  }
  weights {
      name = "myWeightName2"
      weight = 0
      default_value = 5.22
      maximize = false
  }
}