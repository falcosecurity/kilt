package cfnpatcher

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Jeffail/gabs/v2"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9][^\W_]?`)

func getParameterName(envarName string) string {
	// Use the envar name if it does not contain any non-alphanumeric chars
	if found := nonAlphanumericRegex.MatchString(envarName); !found {
		return envarName
	}

	// Otherwise, try to make it more readable, e.g. MY_AWESOME_ENVAR becomes myAwesomeEnvar
	parameterName := nonAlphanumericRegex.ReplaceAllStringFunc(strings.ToLower(envarName), func(str string) string {
		return strings.ToUpper(str[1:])
	})
	return parameterName
}

func getOptTags(template *gabs.Container) map[string]string {
	optTags := make(map[string]string)
	if !template.Exists("Properties", "Tags") {
		return optTags
	}
	for _, tag := range template.S("Properties", "Tags").Children() {
		if tag.Exists("Key") && tag.Exists("Value") {
			k, ok := tag.S("Key").Data().(string)
			if !ok {
				panic(fmt.Errorf("tag has an unsupported key type: %s", tag.String()))
			}
			if isOptTagKey(k) {
				v, ok := tag.S("Value").Data().(string)
				if !ok {
					panic(fmt.Errorf("OptIn/OptOut tag %s has an unsupported value type: %s", k, v))
				}
				optTags[k] = v
			}
		}
	}
	return optTags
}

func exitErrorf(msg string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
