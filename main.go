package main

import (
	"fmt"
	"github.com/quinn/qen/gen"
	"os"
)

type Field struct {
	Name     string
	Type     string
	IsSlice  bool
	Children map[string]*Field
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the template file path as an argument")
		return
	}

	templatePath := os.Args[1]
	err := gen.Exec(templatePath)
	if err != nil {
		fmt.Printf("Error generating structs: %v\n", err)
	}
}
