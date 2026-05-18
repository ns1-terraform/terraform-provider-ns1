resource "ns1_apikey" "apikey" {
  #required
  name = "my api key"

  #optional
  teams = ["myteam"]

  #permissions are available at the top level
}

# Example: API key with automatic secret rotation
resource "ns1_apikey" "rotating_key" {
  name = "rotating-api-key"
  
  # Enable automatic secret rotation every 30 days
  # Valid values: "10d", "30d", "90d"
  expiry_duration = "30d"

  # Configure permissions
  dns_view_zones  = true
  dns_manage_zones = true
}

# The secrets are automatically managed and can be viewed in the state
output "rotating_key_secrets" {
  value = ns1_apikey.rotating_key.secrets
  sensitive = true
}
