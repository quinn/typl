package templates

import (
	"bytes"
	"fmt"
	"text/template"
)

type RootArrayInput struct {
	Root []RootArrayInputRoot
}

type RootArrayInputRoot struct {
	Done bool
	Title string
}

func RootArray(input RootArrayInput) (string, error) {
	tmpl, err := template.ParseFiles("templates/root_array.tpl")
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, input)
	if err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return buf.String(), nil
}
