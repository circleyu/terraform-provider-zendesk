# API reference:
#   https://developer.zendesk.com/api-reference/ticketing/users/users/

resource "zendesk_users" "example" {
  name  = "John Doe"
  email = "john.doe@example.com"
  role  = "end-user"
}

