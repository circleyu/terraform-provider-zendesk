# API reference:
#   https://developer.zendesk.com/api-reference/ticketing/account-configuration/custom_roles/

resource "zendesk_custom_roles" "example" {
  name        = "Example Custom Role"
  description = "An example custom role"
}

