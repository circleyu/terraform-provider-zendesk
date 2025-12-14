# API reference:
#   https://developer.zendesk.com/api-reference/ticketing/tickets/ticket-statuses/

resource "zendesk_custom_statuses" "example" {
  status_category = "open"
  agent_label     = "In Progress"
  end_user_label  = "Working on it"
  active          = true
}

