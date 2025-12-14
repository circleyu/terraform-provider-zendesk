package zendesk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
)

func TestIsValidFile(t *testing.T) {
	v := isValidFile()
	path := cty.GetAttrPath("file_path")

	diags := v("testdata/street.jpg", path)
	if diags.HasError() {
		t.Fatalf("is Valid returned an error")
	}

	diags = v("Missing", path)
	if !diags.HasError() {
		t.Fatalf("is Valid did not return an error for missing file")
	}

	diags = v("testdata", path)
	if !diags.HasError() {
		t.Fatalf("is Valid did not return an error for a directory")
	}
}

func readExampleConfig(t *testing.T, filename string) string {
	dir, err := filepath.Abs("../examples")
	if err != nil {
		t.Fatalf("Failed to resolve fixture directory. Check the path: %s", err)
	}

	filepath := filepath.Join(dir, filename)
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Failed to read fixture. %v", err)
	}

	return string(bytes)
}

func concatExampleConfig(t *testing.T, configs ...string) string {
	builder := new(strings.Builder)
	for _, config := range configs {
		_, err := fmt.Fprintln(builder, config)
		if err != nil {
			t.Fatalf("Encountered an error while concatonating config files: %v", err)
		}
	}

	return builder.String()
}
