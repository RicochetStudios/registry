package registry

import (
	"fmt"
	"regexp"

	polarisv1alpha1 "github.com/RicochetStudios/polaris/apis/v1alpha1"
)

const (
	// templateRegex is a regular expression to validate templates.
	templateRegex string = `^{{ (?P<tpl>(\.\w+)*) }}$`
)

// TemplateValue takes a value and resolves its template if it is a template.
// If the value is not a template, it returns the value as is.
func TemplateValue(v string, s Schema, i polarisv1alpha1.ServerSpec) string {
	// Template the env var if needed.
	re, err := regexp.Compile(templateRegex)
	// If the regex is invalid, return an empty string.
	if err != nil {
		return ""
	}

	if re.MatchString(v) {
		// Get the template to target.
		matches := re.FindStringSubmatch(v)
		tplIndex := re.SubexpIndex("tpl")
		tpl := matches[tplIndex]

		// Resolve the templates.
		switch tpl {
		case ".name":
			return i.Name
		case ".modLoader":
			return i.Game.ModLoader
		case ".players":
			return fmt.Sprint(s.Sizes[i.Size].Players)
		case ".version":
			return i.Game.Version
		}

	}

	// If it is not a template, return an empty string.
	return v
}
