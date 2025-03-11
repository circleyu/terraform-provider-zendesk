resource "zendesk_macro" "temp_utkarsh_test_macro" {
  title = "Macro TF name modified"
  description = "Macro TF description add something"

  action {
    field = "subject"
    value = "foo bar temp terraform"
  }
  action {
    field = "status"
    value = "solved"
  }
}
