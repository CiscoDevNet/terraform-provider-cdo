package acctest

import (
	"bytes"
	"fmt"
	"reflect"
	"text/template"
)

// given a string template and an object,
// render the template with the object and return the rendered string.
func MustParseTemplate(tmpl string, obj any) string {
	buf := bytes.Buffer{}
	if err := template.Must(template.New("").Parse(tmpl)).Execute(&buf, obj); err != nil {
		panic(err)
	}
	return buf.String()
}

func MustOverrideFields[K any](obj K, fields map[string]any) K {
	copyObj := obj
	copyValue := reflect.ValueOf(&copyObj).Elem()

	for k, v := range fields {
		field := copyValue.FieldByName(k)
		if !field.IsValid() || !field.CanSet() {
			panic(fmt.Sprintf("'%s' is an invalid field", k))
		}
		field.Set(reflect.ValueOf(v))
	}

	return copyObj
}
