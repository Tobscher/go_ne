package shared

import "fmt"

// ExtractOptions returns an array of strings from the given interface.
// This is useful if your options are defined as follows:
// plugin: foo
// options:
//   packages:
//     - "foo"
//     - "bar"
//     - "baz"
func ExtractOptions(p interface{}) []string {
	var result []string
	raw := p.([]interface{})

	for _, v := range raw {
		result = append(result, fmt.Sprintf("%v", v))
	}

	return result
}

// ExtractBool converts the given interface to a boolean.
func ExtractBool(p interface{}) bool {
	return p.(bool)
}

// ExtractTruthy converts the given string to a boolean.
func ExtractTruthy(p interface{}) bool {
	switch t := p.(type) {
	case string:
		return t == "true" || t == "yes"
	case bool:
		return t
	}

	return false
}

// ExtractString converts the given interface to a string.
func ExtractString(p interface{}) string {
	return p.(string)
}
