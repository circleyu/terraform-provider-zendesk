package zendesk

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type getter interface {
	Get(string) interface{}
	GetOk(string) (interface{}, bool)
}

type setter interface {
	Set(string, interface{}) error
}

type identifiable interface {
	Id() string
	SetId(string)
}

type identifiableGetterSetter interface {
	identifiable
	getter
	setter
}

type mapGetterSetter map[string]interface{}

func (m mapGetterSetter) Get(k string) interface{} {
	v, ok := m[k]
	if !ok {
		return nil
	}

	return v
}

func (m mapGetterSetter) GetOk(k string) (interface{}, bool) {
	v, ok := m[k]
	return v, ok
}

func (m mapGetterSetter) Set(k string, v interface{}) error {
	m[k] = v
	return nil
}

type identifiableMapGetterSetter struct {
	mapGetterSetter
	id string
}

func newIdentifiableGetterSetter() identifiableGetterSetter {
	return &identifiableMapGetterSetter{
		mapGetterSetter: make(mapGetterSetter),
	}
}

func (i *identifiableMapGetterSetter) Id() string {
	return i.id
}

func (i *identifiableMapGetterSetter) SetId(id string) {
	i.id = id
}

func isValidFile() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		v, ok := i.(string)
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid value type",
				Detail:   fmt.Sprintf("expected type of %s to be string", pathString(path)),
			})
			return diags
		}

		f, err := os.Stat(v)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "File not found",
				Detail:   err.Error(),
			})
			return diags
		}

		if f.IsDir() {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid file path",
				Detail:   fmt.Sprintf("%s: %s is a directory", pathString(path), v),
			})
			return diags
		}

		return diags
	}
}

func pathString(path cty.Path) string {
	if len(path) == 0 {
		return "value"
	}
	// Simple conversion - just use the last step
	if len(path) > 0 {
		lastStep := path[len(path)-1]
		if attrStep, ok := lastStep.(cty.GetAttrStep); ok {
			return attrStep.Name
		}
	}
	return "value"
}

func setSchemaFields(d setter, m map[string]interface{}) error {
	for k, v := range m {
		err := d.Set(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func atoi64(anum string) (int64, error) {
	return strconv.ParseInt(anum, 10, 64)
}

func debugLog(jsonableData interface{}, desc string) error {
	marshaled, err := json.MarshalIndent(jsonableData, "", "   ")
	if err != nil {
		return err
	}
	fmt.Printf("### LOG %s START ###\n", desc)
	fmt.Println(string(marshaled))
	fmt.Printf("### LOG %s END ###\n", desc)
	return nil
}
