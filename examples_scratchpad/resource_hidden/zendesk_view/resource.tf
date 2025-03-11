resource "zendesk_view" "temputkarsh-tf" {
  title = "TEMP TF VIEW"
  description = "This is by terraform"
  all {
    field    = "status"
    operator = "is"
    value    = "pending"
  }
  any {
    field    = "status"
    operator = "is"
    value    = "open"
  }
  sort_by = "requester"
  group_by = "status"
  group_order = "asc"
  sort_order = "asc"
  columns = ["subject", "status", "18429918055186"]
  restrictions = [18373407148562]
}