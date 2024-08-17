## typl: generate structs for go templates

```
go install github.com/quinn/typl@latest
```

```
# templates/todo_list.tpl

<h1>{{.PageTitle}}</h1>
<ul>
    {{range .Todos}}
        {{if .Done}}
            <li class="done">{{.Title}}</li>
        {{else}}
            <li>{{.Title}}</li>
        {{end}}
    {{end}}
</ul>
```

run:

```
typl templates/todo_list.tpl
```

creates:

```go
// templates/todo_list.go

package templates

import (
	"bytes"
	"fmt"
	"text/template"
)

type TodoListInput struct {
	PageTitle string
	Todos []TodoListInputTodos
}

type TodoListInputTodos struct {
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
```
