package zendesk

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nukosuke/go-zendesk/zendesk"
	"github.com/nukosuke/terraform-provider-zendesk/zendesk/client"
	"github.com/nukosuke/terraform-provider-zendesk/zendesk/models"
)

// mockTicketFormAPI is a mock implementation of client.TicketFormAPI
type mockTicketFormAPI struct {
	createTicketForm func(ctx context.Context, ticketForm models.TicketForm) (models.TicketForm, error)
	getTicketForm    func(ctx context.Context, id int64) (models.TicketForm, error)
	deleteTicketForm func(ctx context.Context, id int64) error
	updateTicketForm func(ctx context.Context, id int64, form models.TicketForm) (models.TicketForm, error)
}

func (m *mockTicketFormAPI) CreateTicketForm(ctx context.Context, ticketForm models.TicketForm) (models.TicketForm, error) {
	if m.createTicketForm != nil {
		return m.createTicketForm(ctx, ticketForm)
	}
	return models.TicketForm{}, nil
}

func (m *mockTicketFormAPI) GetTicketForm(ctx context.Context, id int64) (models.TicketForm, error) {
	if m.getTicketForm != nil {
		return m.getTicketForm(ctx, id)
	}
	return models.TicketForm{}, nil
}

func (m *mockTicketFormAPI) DeleteTicketForm(ctx context.Context, id int64) error {
	if m.deleteTicketForm != nil {
		return m.deleteTicketForm(ctx, id)
	}
	return nil
}

func (m *mockTicketFormAPI) UpdateTicketForm(ctx context.Context, id int64, form models.TicketForm) (models.TicketForm, error) {
	if m.updateTicketForm != nil {
		return m.updateTicketForm(ctx, id, form)
	}
	return models.TicketForm{}, nil
}

func (m *mockTicketFormAPI) GetTicketForms(ctx context.Context, options *zendesk.TicketFormListOptions) ([]models.TicketForm, zendesk.Page, error) {
	return nil, zendesk.Page{}, nil
}

func TestCreateTicketForm(t *testing.T) {
	i := newIdentifiableGetterSetter()
	out := models.TicketForm{
		ID:   12345,
		Name: "foo",
	}

	m := &mockTicketFormAPI{
		createTicketForm: func(ctx context.Context, ticketForm models.TicketForm) (models.TicketForm, error) {
			return out, nil
		},
	}

	if diags := createTicketForm(context.Background(), i, m); len(diags) != 0 {
		t.Fatal("create ticket field returned an error")
	}

	if v := i.Id(); v != "12345" {
		t.Fatalf("Create did not set resource id. Id was %s", v)
	}

	if v := i.Get("name"); v != "foo" {
		t.Fatalf("Create did not set resource name. name was %v", v)
	}
}

func TestDeleteTicketForm(t *testing.T) {
	i := newIdentifiableGetterSetter()
	i.SetId("12345")

	m := &mockTicketFormAPI{
		deleteTicketForm: func(ctx context.Context, id int64) error {
			if id != 12345 {
				t.Fatalf("Expected ID 12345, got %d", id)
			}
			return nil
		},
	}

	if diags := deleteTicketForm(context.Background(), i, m); len(diags) != 0 {
		t.Fatal("create ticket field returned an error")
	}
}

func TestReadTicketForm(t *testing.T) {
	i := newIdentifiableGetterSetter()
	i.SetId("12345")

	expected := models.TicketForm{
		Name:     "foobar",
		Position: int64(1),
	}

	m := &mockTicketFormAPI{
		getTicketForm: func(ctx context.Context, id int64) (models.TicketForm, error) {
			if id != 12345 {
				t.Fatalf("Expected ID 12345, got %d", id)
			}
			return expected, nil
		},
	}

	if diags := readTicketForm(context.Background(), i, m); len(diags) != 0 {
		t.Fatalf("recieved an error when calling read ticket form: %v", diags)
	}
}

func TestUnmarshalTicketForm(t *testing.T) {

	d := &identifiableMapGetterSetter{
		id: "47",
		mapGetterSetter: mapGetterSetter{
			"url":              "https://company.zendesk.com/api/v2/ticket_forms/47.json",
			"name":             "Snowboard Problem",
			"display_name":     "Snowboard Damage",
			"end_user_visible": true,
			"position":         9999,
			"active":           true,
			"default":          false,
			"in_all_brands":    false,
		},
	}

	tf, err := unmarshalTicketForm(d)
	if err != nil {
		t.Fatalf("unmarshal returned an error: %v", err)
	}

	if tf.Name != d.Get("name") {
		t.Fatalf("ticket did not have the correct name")
	}
}

func testTicketFormDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(client.TicketFormAPI)

	for k, r := range s.RootModule().Resources {
		if r.Type != "zendesk_ticket_form" {
			continue
		}

		id, err := atoi64(r.Primary.ID)
		if err != nil {
			return err
		}

		form, err := client.GetTicketForm(context.Background(), id)
		if err != nil {
			return fmt.Errorf("got an error from zendesk when trying to fetch the destroyed ticket form %s. %v", k, err)
		}

		if form.Active {
			return fmt.Errorf("form %v is still active", form)
		}

	}

	return nil
}

func TestAccTicketFormExample(t *testing.T) {
	configs := []string{
		readExampleConfig(t, "resources/zendesk_ticket_field/resource.tf"),
		readExampleConfig(t, "resources/zendesk_ticket_form/resource.tf"),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testTicketFieldDestroyed,
			testTicketFormDestroyed,
		),
		Steps: []resource.TestStep{
			{
				Config: concatExampleConfig(t, configs...),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zendesk_ticket_form.form-1", "name", "Form 1"),
					resource.TestCheckResourceAttr("zendesk_ticket_form.form-2", "name", "Form 2"),
				),
			},
		},
	})
}
