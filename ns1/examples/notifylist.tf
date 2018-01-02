resource "ns1_notifylist" "test" {
  #required
  name = "terraform test"

  #optional
  notifications = {
    type = "webhook"
    config = {
      url = "http://localhost:9090"
    }
  }
}