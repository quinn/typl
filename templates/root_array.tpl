{{range .}}
    {{if .Done}}
        <li class="done">{{.Title}}</li>
    {{else}}
        <li>{{.Title}}</li>
    {{end}}
{{end}}
