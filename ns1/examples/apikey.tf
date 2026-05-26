resource "ns1_apikey" "apikey" {
  #required
  name = "my api key"

  #optional
  teams = ["myteam"]

  #permissions are available at the top level
}

# Example: API key with secret expiration
resource "ns1_apikey" "expiring_key" {
  name = "expiring-api-key"
  
  # Set secret expiration period to 30 days
  # Secrets will expire after this period and must be manually renewed
  # Accepts any duration in '<number>d' format (e.g., "10d", "30d", "90d")
  expiry_duration = "30d"

  # Configure permissions
  dns_view_zones  = true
  dns_manage_zones = true
}

# The secrets metadata can be viewed in the state
output "expiring_key_secrets" {
  value = ns1_apikey.expiring_key.secrets
  sensitive = true
}
