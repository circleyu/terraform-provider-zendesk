
resource "zendesk_organization_field" "foobar" {
  title       = "field"
  type        = "checkbox"
  key         = "foobar"
  description = "foo bar some desc"
}

resource "zendesk_dynamic_content" "foodcnew" {
  name = "dcutkok"
}

resource "zendesk_dynamic_content" "life_isgood" {
  name = "lifeisgoodright"
}
resource "zendesk_dynamic_content_variant" "foobar" {
  content                 = "FooBar change is good"
  locale_id               = 2
  default                 = true
  dynamic_content_item_id = zendesk_dynamic_content.life_isgood.id
}

resource "zendesk_dynamic_content" "better" {
  name = "better"
}

resource "zendesk_dynamic_content_variant" "better_content" {
  content                 = "Some data here is template good"
  locale_id               = 1
  default                 = true
  dynamic_content_item_id = 18847895265170
}

resource "zendesk_dynamic_content_variant" "foobaranother" {
  content                 = "foo bar is here right??? new"
  locale_id               = 1
  default                 = true
  dynamic_content_item_id = zendesk_dynamic_content.life_isgood.id
}

resource "zendesk_macro" "temp_utkarsh_test_macro" {
  title       = "NEW: Macro TF name modified"
  description = "Macro TF description add something"

  action {
    field = "subject"
    value = "foo bar temp terraform"
  }
  action {
    field = "status"
    value = "solved"
  }

  restrictions = [18651233156242, 17784599337874]
}

resource "zendesk_ticket_field" "checkbox-field3" {
  title                 = "Checkbox Field (this is a test)"
  type                  = "regexp"
  regexp_for_validation = "^foobar$"
  description           = "test here something"
  agent_description     = "agent description change here"
  required              = true
}

resource "zendesk_user_field" "tagger-field" {
  title = "Tagger Field"
  type  = "dropdown"
  key   = "foobar"

  custom_field_option {
    id    = 20907892899730
    name  = "option b here foo bar"
    value = "optb"
  }
  custom_field_option {
    id    = 20909205605650
    name  = "option a"
    value = "opta"
  }


}


resource "zendesk_webhook" "example-basic-auth-webhook" {
  name           = "Example Webhook with Basic Auth"
  endpoint       = "https://example.com/status/200"
  http_method    = "POST"
  request_format = "form_encoded"
  status         = "active"
  subscriptions  = ["conditional_ticket_events"]

  authentication {
    type         = "basic_auth"
    add_position = "header"
    data = {
      username = "john.doe"
      password = "password+change"
    }
  }
}

# {
#                 "parent_field_id": 360000035929,
#                 "parent_field_type": "checkbox",
#                 "value": true,
#                 "child_fields": [
#                     {
#                         "id": 33735105,
#                         "is_required": false,
#                         "required_on_statuses": {
#                             "type": "NO_STATUSES"
#                         }
#                     }
#                 ]
#             },


# resource "zendesk_ticket_form" "form-1" {
#   name = "Form 1"
#   ticket_field_ids = [
#     20713391817234,
#     20713391818514,
#     20713412708626,
#     20713391817490,
#     20713391816338,
#     20713412707218,
#     20975060928402,
#     20975076773266,
#     20975044871186,
#     20975076776594,
#     20975029692434,
#     20975013788818, //zeit
#     20975043203218,
#   ]
# }

