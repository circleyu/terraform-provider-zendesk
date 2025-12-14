# API reference:
#   https://developer.zendesk.com/api-reference/ticketing/queues/

resource "zendesk_queues" "example" {
  name        = "Example Queue"
  description = "An example queue"
}

