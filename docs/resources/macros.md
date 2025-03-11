```terraform

resource "zendesk_macro" "test_macro" {
    title = "Macro Title"
    description = "Macro description add something"
    action {
        field = "subject"
        value = "subject here"
    }
    action {
        field = "status"
        value = "solved"
    }
    action {
        field = "side_conversation"
        value = jsonencode(["this is the subject", "<p>this is the email body</p>", "text/html"])
    }
}
```
