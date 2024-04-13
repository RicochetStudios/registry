package registry

import (
	"testing"

	polarisv1alpha1 "github.com/RicochetStudios/polaris/apis/v1alpha1"
	"github.com/google/go-cmp/cmp"
)

// exampleServerSpec defines a minimal ServerSpec for testing.
var exampleServerSpec = polarisv1alpha1.ServerSpec{
	Name: "minecraft",
	Size: "s",
	Game: polarisv1alpha1.Game{
		Version:   "1.17.1",
		ModLoader: "forge",
	},
}

// TestTemplate tests that TemplateValue() will return the correct value.
func TestTemplate(t *testing.T) {
	// Test matrix.
	var matrix = []struct {
		name  string
		input string
		want  string
	}{
		{"correct template: .name", "{{ .name }}", "minecraft"},
		{"correct template: .modLoader", "{{ .modLoader }}", "forge"},
		{"correct template: .players", "{{ .players }}", "16"},
		{"correct template: .version", "{{ .version }}", "1.17.1"},
		{"no template", "EULA", "EULA"},
		{"broken template", "{{ .players", "{{ .players"},
	}

	// Run tests.
	for _, test := range matrix {
		// Call the function to test.
		got := TemplateValue(test.input, exampleSchema, exampleServerSpec)

		// Error if results are incorrect.
		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Fatalf(`TemplateValue() test case "%s" mistmatch (-want +got):\n%s`, test.name, diff)
		}
	}
}
