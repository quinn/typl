package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"text/template/parse"
	"unicode"
)

type Field struct {
	Name     string
	Type     string
	IsSlice  bool
	Children map[string]*Field
}

// Exec is the main function that can be called with a template path, output path, and package name
func Exec(templatePath, outputPath, packageName string) error {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	base := strings.TrimSuffix(filepath.Base(templatePath), filepath.Ext(templatePath))
	funcName := toCamelCase(base)
	structName := funcName + "Input"

	fields, err := extractFieldsFromTemplateAST(string(content))
	if err != nil {
		return fmt.Errorf("error extracting fields: %v", err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	fmt.Fprintf(outputFile, "package %s\n\n", packageName)
	fmt.Fprintf(outputFile, "import (\n\t\"bytes\"\n\t\"fmt\"\n\t\"text/template\"\n)\n\n")

	generateStructs(outputFile, structName, fields)
	generateRenderFunction(outputFile, funcName, structName, templatePath)
	return nil
}

// Update other functions to write to the provided file instead of stdout
func generateStructs(w *os.File, structName string, fields map[string]*Field) {
	generateStruct(w, structName, fields, 0, structName)
}

func generateStruct(w *os.File, structName string, fields map[string]*Field, indent int, prefix string) {
	indentStr := strings.Repeat("\t", indent)
	fmt.Fprintf(w, "%stype %s struct {\n", indentStr, structName)
	for _, field := range fields {
		fieldType := field.Type
		if field.IsSlice {
			fieldType = "[]" + prefix + strings.TrimSuffix(field.Type, "Item")
		}
		fmt.Fprintf(w, "%s\t%s %s\n", indentStr, strings.Title(field.Name), fieldType)
	}
	fmt.Fprintf(w, "%s}\n\n", indentStr)

	for _, field := range fields {
		if len(field.Children) > 0 {
			childStructName := prefix + strings.TrimSuffix(field.Type, "Item")
			if !field.IsSlice {
				childStructName = prefix + field.Type
			}
			generateStruct(w, childStructName, field.Children, indent, prefix)
		}
	}
}

func generateRenderFunction(w *os.File, funcName, structName, templatePath string) {
	fmt.Fprintf(w, "func %s(input %s) (string, error) {\n", funcName, structName)
	fmt.Fprintf(w, "\ttmpl, err := template.ParseFiles(%q)\n", templatePath)
	fmt.Fprintf(w, "\tif err != nil {\n")
	fmt.Fprintf(w, "\t\treturn \"\", fmt.Errorf(\"error parsing template: %%v\", err)\n")
	fmt.Fprintf(w, "\t}\n\n")
	fmt.Fprintf(w, "\tvar buf bytes.Buffer\n")
	fmt.Fprintf(w, "\terr = tmpl.Execute(&buf, input)\n")
	fmt.Fprintf(w, "\tif err != nil {\n")
	fmt.Fprintf(w, "\t\treturn \"\", fmt.Errorf(\"error executing template: %%v\", err)\n")
	fmt.Fprintf(w, "\t}\n\n")
	fmt.Fprintf(w, "\treturn buf.String(), nil\n")
	fmt.Fprintf(w, "}\n")
}

func toCamelCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, "")
}

func extractFieldsFromTemplateAST(content string) (map[string]*Field, error) {
	fields := make(map[string]*Field)

	tmpl, err := template.New("template").Parse(content)
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %v", err)
	}

	var extractFields func(node parse.Node, currentFields map[string]*Field) error
	extractFields = func(node parse.Node, currentFields map[string]*Field) error {
		if node == nil {
			return nil
		}

		switch n := node.(type) {
		case *parse.ActionNode:
			if n.Pipe == nil {
				return nil
			}
			for _, cmd := range n.Pipe.Cmds {
				for _, arg := range cmd.Args {
					if field, ok := arg.(*parse.FieldNode); ok && len(field.Ident) > 0 {
						addField(currentFields, field.Ident, false)
					}
				}
			}
		case *parse.ListNode:
			if n != nil {
				for _, item := range n.Nodes {
					if err := extractFields(item, currentFields); err != nil {
						return err
					}
				}
			}
		case *parse.IfNode:
			if n.Pipe != nil {
				for _, cmd := range n.Pipe.Cmds {
					for _, arg := range cmd.Args {
						if field, ok := arg.(*parse.FieldNode); ok && len(field.Ident) > 0 {
							addField(currentFields, field.Ident, true)
						}
					}
				}
			}
			if err := extractFields(n.List, currentFields); err != nil {
				return err
			}
			if err := extractFields(n.ElseList, currentFields); err != nil {
				return err
			}
		case *parse.RangeNode:
			if n.Pipe != nil && len(n.Pipe.Decl) > 0 {
				rangeVar := n.Pipe.Decl[0].String()
				if _, exists := currentFields[rangeVar]; !exists {
					currentFields[rangeVar] = &Field{
						Name:     rangeVar,
						Type:     rangeVar + "Item",
						IsSlice:  true,
						Children: make(map[string]*Field),
					}
				}
				if err := extractFields(n.List, currentFields[rangeVar].Children); err != nil {
					return err
				}
			} else if n.Pipe != nil {
				for _, cmd := range n.Pipe.Cmds {
					for _, arg := range cmd.Args {
						if field, ok := arg.(*parse.FieldNode); ok && len(field.Ident) > 0 {
							rangeVar := field.Ident[0]
							if _, exists := currentFields[rangeVar]; !exists {
								currentFields[rangeVar] = &Field{
									Name:     rangeVar,
									Type:     rangeVar + "Item",
									IsSlice:  true,
									Children: make(map[string]*Field),
								}
							}
							if err := extractFields(n.List, currentFields[rangeVar].Children); err != nil {
								return err
							}
						}
					}
				}
			}
			if err := extractFields(n.ElseList, currentFields); err != nil {
				return err
			}
		case *parse.WithNode:
			if err := extractFields(n.List, currentFields); err != nil {
				return err
			}
			if err := extractFields(n.ElseList, currentFields); err != nil {
				return err
			}
		}
		return nil
	}

	if err := extractFields(tmpl.Tree.Root, fields); err != nil {
		return nil, fmt.Errorf("error extracting fields: %v", err)
	}

	return fields, nil
}

func addField(fields map[string]*Field, ident []string, isConditional bool) {
	if len(ident) == 0 {
		return
	}

	fieldName := ident[0]
	if _, exists := fields[fieldName]; !exists {
		fieldType := "string"
		if isConditional {
			fieldType = "bool"
		}
		fields[fieldName] = &Field{
			Name:     fieldName,
			Type:     fieldType,
			Children: make(map[string]*Field),
		}
	}

	if len(ident) > 1 {
		addField(fields[fieldName].Children, ident[1:], isConditional)
	}
}
