package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"text/template/parse"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the template file path as an argument")
		return
	}

	templatePath := os.Args[1]
	content, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	structName := getStructNameFromFileName(templatePath)
	fields := extractFieldsFromTemplateAST(string(content))

	generateStruct(structName, fields)
}

func getStructNameFromFileName(filePath string) string {
	base := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	return strings.Title(base) + "Input"
}

func extractFieldsFromTemplateAST(content string) map[string]string {
	fields := make(map[string]string)

	tmpl, err := template.New("template").Parse(content)
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return fields
	}

	var extractFields func(node parse.Node)
	extractFields = func(node parse.Node) {
		switch n := node.(type) {
		case *parse.ActionNode:
			for _, cmd := range n.Pipe.Cmds {
				for _, arg := range cmd.Args {
					if field, ok := arg.(*parse.FieldNode); ok {
						fieldName := field.Ident[0]
						if _, exists := fields[fieldName]; !exists {
							fields[fieldName] = "string" // Default to string, you can enhance this later
						}
					}
				}
			}
		case *parse.ListNode:
			for _, item := range n.Nodes {
				extractFields(item)
			}
		case *parse.IfNode:
			extractFields(n.List)
			extractFields(n.ElseList)
		case *parse.RangeNode:
			extractFields(n.List)
			extractFields(n.ElseList)
		case *parse.WithNode:
			extractFields(n.List)
			extractFields(n.ElseList)
		}
	}

	extractFields(tmpl.Tree.Root)
	return fields
}

func generateStruct(structName string, fields map[string]string) {
	fmt.Printf("type %s struct {\n", structName)
	for fieldName, fieldType := range fields {
		fmt.Printf("\t%s %s\n", strings.Title(fieldName), fieldType)
	}
	fmt.Println("}")
}
