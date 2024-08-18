package templates

import (
	"bytes"
	"fmt"
	"text/template"
)

type TodoListInput struct {
	PageTitle string
	Todos []TodoListTodos
}

type TodoListTodos struct {
	Done bool
	Title string
}

func TodoList(input TodoListInput) (string, error) {
	tmpl, err := template.ParseFiles("templates/todo_list.tpl")
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
