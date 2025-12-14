# API reference:
#   https://developer.zendesk.com/api-reference/ticketing/tickets/tickets/

resource "zendesk_tickets" "example" {
  subject     = "Example Ticket"
  description = "This is an example ticket"
  priority    = "normal"
  status      = "new"
}

