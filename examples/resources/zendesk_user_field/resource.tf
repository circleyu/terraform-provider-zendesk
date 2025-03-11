# # # API reference:
# # #   https://developer.zendesk.com/rest_api/docs/support/user_fields

# # resource "zendesk_user_field" "checkbox-field" {
# #   title = "Checkbox Field Name Change"
# #   type = "checkbox"
# #   key = "temp_test_utkarsh_checkbox_field"
# #   description = "Check box field has some description"
# # }
# # # resource "zendesk_trigger_category" "test" {
# # #   name = "TestCatWhy"
# # #   position = 11
# # # }

# # # resource "zendesk_trigger" "auto-reply-trigger" {
# # #   title  = "UTK TEST"
# # #   active = true

# # #   all {
# # #     field    = "role"
# # #     operator = "is"
# # #     value    = "end_user"
# # #   }

# # #   all {
# # #     field    = "update_type"
# # #     operator = "is"
# # #     value    = "Create"
# # #   }

# # #   all {
# # #     field    = "status"
# # #     operator = "is_not"
# # #     value    = "solved"
# # #   }

# # #   action {
# # #     field = "notification_user"
# # #     value = jsonencode([
# # #       "requester_id",
# # #       "Dear my customer",
# # #       "Hi. This message was configured by terraform-provider-zendesk."
# # #     ])
# # #   }
# # #   category_id = zendesk_trigger_category.test.id
# # # }

# # resource "zendesk_user_field" "date-field" {
# #   title = "Date Field"
# #   type = "date"
# #   key = "temp_test_utkarsh_date_field"
# # }

# # resource "zendesk_user_field" "decimal-field" {
# #   title = "Decimal Field"
# #   type = "decimal"
# #   key = "temp_test_utkarsh_decimal_field"
# # }

# # resource "zendesk_user_field" "integer-field" {
# #   title = "Integer Field"
# #   type = "integer"
# #   key = "temp_test_utkarsh_integer_field"
# # }

# # resource "zendesk_user_field" "regexp-field" {
# #   title = "Regexp Field"
# #   type = "regexp"
# #   regexp_for_validation = "^[0-9]+-[0-9]+-[0-9]+$"
# #   key = "temp_test_utkarsh_regexp_field"
# # }

# # # resource "zendesk_user_field" "tagger-field" {
# # #   title = "Tagger Field"
# # #   type = "tagger"

# # #   custom_field_option {
# # #     name  = "Option 1"
# # #     value = "opt1"
# # #   }

# # #   custom_field_option {
# # #     name  = "Option 2"
# # #     value = "opt2"
# # #   }
# # # }

# # resource "zendesk_user_field" "text-field" {
# #   title = "Text Field"
# #   type = "text"
# #   key = "temp_test_utkarsh_text_field"
# # }

# # resource "zendesk_user_field" "textarea-field" {
# #   title = "Textarea Field"
# #   type = "textarea"
# #   key = "temp_test_utkarsh_textarea_field"
# # }

# # # data "zendesk_user_field" "assignee" {
# # #   type = "assignee"
# # # }

# # # data "zendesk_user_field" "group" {
# # #   type = "group"
# # # }

# # # data "zendesk_user_field" "status" {
# # #   type = "status"
# # # }

# # # data "zendesk_user_field" "subject" {
# # #   type = "subject"
# # # }

# # # data "zendesk_user_field" "description" {
# # #   type = "description"
# # # }

# # resource "zendesk_view" "temputkarsh-tfA" {
# #   title = "TEMP TF VIEW A"
# #   description = "This is by terraform"
# #   position = 9
# #   all {
# #     field    = "status"
# #     operator = "is"
# #     value    = "pending"
# #   }
# #   any {
# #     field    = "status"
# #     operator = "is"
# #     value    = "open"
# #   }
# #   sort_by = "requester"
# #   group_by = "status"
# #   group_order = "asc"
# #   sort_order = "asc"
# #   columns = ["status", "subject"]
# # }

# resource "zendesk_view" "temputkarsh-tfB" {
#   title = "TEMP TF VIEW B"
#   description = "This is by terraform"
#   position = 9
#   all {
#     field    = "status"
#     operator = "is"
#     value    = "pending"
#   }
#   any {
#     field    = "status"
#     operator = "is"
#     value    = "open"
#   }
#   sort_by = "requester"
#   group_by = "status"
#   group_order = "asc"
#   sort_order = "asc"
#   columns = ["status", "subject"]
#   restrictions = [18373407148562]
# }
