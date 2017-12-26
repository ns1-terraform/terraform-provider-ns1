resource "ns1_apikey" "apikey" {
  #required
  name = "my api key"

  #optional
  teams = ["myteam"]
}